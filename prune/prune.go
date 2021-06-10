package prune

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/chr4/pwgen"
	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/generate"
	"github.com/cycloidio/inframap/provider/factory"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/states/statefile"
	uuid "github.com/satori/go.uuid"
)

// reARN matches an arn string
var reARN = regexp.MustCompile("^arn:*")

// Prune will prune the tfstate of unneeded information and if replaceCanonicals is specified
// the resource canonicals will also be changed, for exmple 'aws_lb.front' will be changed to
// a random name like 'aws_lb.XptaK'
func Prune(tfstate json.RawMessage, replaceCanonicals bool) (json.RawMessage, error) {
	// Since TF 0.13 'depends_on' has been dropped, so we do a manual
	// replace from '"depends_on"' to '"dependencies"'
	tfstate = bytes.ReplaceAll(tfstate, []byte("\"depends_on\""), []byte("\"dependencies\""))
	err := generate.ValidateTFStateVersion(tfstate)
	if err != nil {
		return nil, fmt.Errorf("error while validating TFStateVersion: %w", err)
	}

	buf := bytes.NewBuffer(tfstate)

	file, err := statefile.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error while reading TFState: %w", err)
	}

	// canonicals holds the old canonical of the resource "aws_elb.front" as key
	// and as value the new name it has been randomly given
	canonicals := make(map[string]string)
	for _, m := range file.State.Modules {
		removeKeys := make([]string, 0)
		for rk, rv := range m.Resources {
			// If it's not a Resource we ignore it
			if rv.Addr.Resource.Mode != addrs.ManagedResourceMode {
				removeKeys = append(removeKeys, rk)
				continue
			}

			pv, rs, err := factory.GetProviderAndResource(rk)
			if err != nil {
				if errors.Is(err, errcode.ErrProviderNotFound) {
					removeKeys = append(removeKeys, rk)
					continue
				}
				return nil, err
			}

			// If it's not a Node or Edge we append it to delete
			if !pv.IsNode(rs) && !pv.IsEdge(rs) {
				removeKeys = append(removeKeys, rk)
				continue
			}

			attrs := pv.UsedAttributes()
			for _, iv := range rv.Instances {
				if replaceCanonicals {
					canonicals[rk] = fmt.Sprintf("%s.%s", rs, pwgen.Alpha(5))
				}
				aux := make(map[string]interface{})
				var legacy bool
				if iv.Current.AttrsJSON != nil {
					// For TF 0.12
					err = json.Unmarshal(iv.Current.AttrsJSON, &aux)
					if err != nil {
						return nil, fmt.Errorf("invalid fomrat JSON for resource %q with AttrsJSON %s: %w", string(iv.Current.AttrsJSON), rk, err)
					}
				} else {
					// For TF 0.11
					// AttrsFlat it's a map[string]string so it has to be converted
					// to map[string]interface{} to fit on the aux definition
					legacy = true
					for k, v := range iv.Current.AttrsFlat {
						aux[k] = v
					}
				}

				// Remove the "private" info as we do not need it
				iv.Current.Private = []byte{}

				// TODO: Think on a more provider agnostic solution for this
				if v, ok := aux["id"]; ok {
					vs, ok := v.(string)
					if ok {
						if reARN.MatchString(vs) {
							aux["id"] = uuid.NewV4().String()
						}
					}
				}

				for k := range aux {
					var found bool
					for _, a := range attrs {
						if k == a || regexp.MustCompile(fmt.Sprintf(`^%s\.`, a)).MatchString(k) {
							// One match on the whitelist and we
							// have to break the loop
							found = true
						}
					}

					if !found {
						delete(aux, k)
					}
				}
				if !legacy {
					b, err := json.Marshal(aux)
					if err != nil {
						return nil, err
					}
					iv.Current.AttrsJSON = b
				} else {
					// We need to empty the AttrsFlat first
					iv.Current.AttrsFlat = make(map[string]string)
					for k, v := range aux {
						iv.Current.AttrsFlat[k] = v.(string)
					}
				}
			}
		}

		// Delete all the resources we do not deal with
		for _, k := range removeKeys {
			delete(m.Resources, k)
		}
	}

	// Now that the actual State is pruned of unneeded data
	// we iterate again to change the canonicals and 'dependencies'
	// if needed
	if replaceCanonicals {
		for _, m := range file.State.Modules {
			for _, rv := range m.Resources {
				if newCan, ok := canonicals[fmt.Sprintf("%s.%s", rv.Addr.Resource.Type, rv.Addr.Resource.Name)]; ok {
					splitCan := strings.Split(newCan, ".")
					rv.Addr.Resource.Type = splitCan[0]
					rv.Addr.Resource.Name = splitCan[1]
				}
				for _, iv := range rv.Instances {
					removeDepends := make([]int, 0)
					deps := make(map[string]struct{})
					for i, d := range iv.Current.Dependencies {
						if newCan, ok := canonicals[fmt.Sprintf("%s.%s", d.Resource.Type, d.Resource.Name)]; ok {
							// If the dependency it's already present
							// do not add repeated ones
							if _, ok := deps[newCan]; ok {
								removeDepends = append(removeDepends, i)
							}
							splitCan := strings.Split(newCan, ".")
							d.Resource.Type = splitCan[0]
							d.Resource.Name = splitCan[1]
							deps[newCan] = struct{}{}
						} else {
							removeDepends = append(removeDepends, i)
						}
						iv.Current.Dependencies[i] = d
					}
					for i, idx := range removeDepends {
						iv.Current.Dependencies = append(iv.Current.Dependencies[:(idx-i)], iv.Current.Dependencies[(idx-i)+1:]...)
					}
				}
			}
		}
	}

	b := bytes.Buffer{}
	statefile.Write(file, &b)

	return json.RawMessage(b.Bytes()), nil
}

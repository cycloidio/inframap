package infraview

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/cycloidio/infraview/errcode"
	"github.com/cycloidio/infraview/factory"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/states/statefile"
)

func Prune(tfstate json.RawMessage) (json.RawMessage, error) {
	buf := bytes.NewBuffer(tfstate)
	file, err := statefile.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error while reading TFState: %w", err)
	}

	for _, m := range file.State.Modules {
		removeKeys := make([]string, 0)
		for rk, rv := range m.Resources {
			// If it's not a Resource we ignore it
			if rv.Addr.Mode != addrs.ManagedResourceMode {
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

				for k, _ := range aux {
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

	b := bytes.Buffer{}
	statefile.Write(file, &b)

	return json.RawMessage(b.Bytes()), nil
}

package main

// objectclasses are generated by the parser from the schema source
type objectclass struct {
	Name                            string
	Desc                            string
	Obsolete                        bool
	Sup                             []string
	Abstract, Structural, Auxiliary bool
	Must                            []string
	May                             []string
}

// create a new objectclass initialized with the values from the map m
func newObjectclass(m map[int]interface{}) *objectclass {
	o := objectclass{}
	for k, v := range m {
		switch k {
		case NAME:
			o.Name = v.(string)
		case DESC:
			o.Desc = v.(string)
		case OBSOLETE:
			o.Obsolete = v.(bool)
		case SUP:
			o.Sup = v.([]string)
		case ABSTRACT:
			o.Abstract = v.(bool)
		case STRUCTURAL:
			o.Structural = v.(bool)
		case AUXILIARY:
			o.Auxiliary = v.(bool)
		case MUST:
			/*
				attrNames := v.([]string)
				o.Must = make([]*attributetype, len(attrNames))
				for i, a := range attrNames {
					o.Must[i] = attributetypedefs[a]
				}
			*/
			o.Must = v.([]string)
		case MAY:
			/*
				attrNames := v.([]string)
				o.May = make([]*attributetype, len(attrNames))
				for i, a := range attrNames {
					o.May[i] = attributetypedefs[a]
				}
			*/
			o.May = v.([]string)
		}
	}
	return &o
}
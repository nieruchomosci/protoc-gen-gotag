package module

import (
	"github.com/nieruchomosci/protoc-gen-gotag/tagger"
	"github.com/fatih/structtag"
	pgs "github.com/lyft/protoc-gen-star"
)

type tagExtractor struct {
	pgs.Visitor
	pgs.DebuggerCommon

	tags map[string]map[string]*structtag.Tags
}

func newTagExtractor(d pgs.DebuggerCommon) *tagExtractor {
	v := &tagExtractor{DebuggerCommon: d}
	v.Visitor = pgs.PassThroughVisitor(v)
	return v
}

func (v *tagExtractor) VisitOneOf(o pgs.OneOf) (pgs.Visitor, error) {
	var tval string
	ok, err := o.Extension(tagger.E_OneofTags, &tval)
	if err != nil {
		return nil, err
	}

	if !ok {
		return v, nil
	}

	tags, err := structtag.Parse(tval)
	if err != nil {
		return nil, err
	}

	msgName := o.Message().Name().UpperCamelCase().String()

	if v.tags[msgName] == nil {
		v.tags[msgName] = map[string]*structtag.Tags{}
	}

	v.tags[msgName][o.Name().UpperCamelCase().String()] = tags

	return v, nil
}

func (v *tagExtractor) VisitField(f pgs.Field) (pgs.Visitor, error) {
	var tval string
	ok, err := f.Extension(tagger.E_Tags, &tval)
	if err != nil {
		return nil, err
	}

	msgName := f.Message().Name().UpperCamelCase().String()

	if f.InOneOf() {
		msgName = f.Message().Name().UpperCamelCase().String() + "_" + f.Name().UpperCamelCase().String()
	}

	if v.tags[msgName] == nil {
		v.tags[msgName] = map[string]*structtag.Tags{}
	}

	if ok {
		tags, err := structtag.Parse(tval)
		v.CheckErr(err)

		v.tags[msgName][f.Name().UpperCamelCase().String()] = tags
	}

	return v, nil
}

func (v *tagExtractor) Extract(f pgs.File) StructTags {
	v.tags = map[string]map[string]*structtag.Tags{}

	v.CheckErr(pgs.Walk(v, f))

	return v.tags
}

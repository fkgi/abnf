# abnf
ABNF LL-parser by golang

## How to use
### Make Parser function
Make parser function that implement Rule interface.
Functions that returns Basic Operator Rule are provided by this package.

For example,
* `abnf.C(Rule...)` function returns Rule that concatenate all argument Rules.
* `abnf.V(rune)` function returns Rule that match argument rune.
* `abnf.ALPHA()` function returns Rule that match any alphanumeric text.
* `abnf.ETX()` function returns Rule that match no more text is left.

Combile multiple Basic Operator by set each function to arguement.
`abnf.C(abnf.V('-'), abnf.ALPHANUM())` will make Rule that match aplhanumeric text taht start with -.

### Parse Text by created Parser
Call `abnf.ParseString(string, Rule)` function with created Rule. The string will parsed and parsed item Tree is returned.

If text is not match with Rule, nil is returned.

### Get parsed value from Tree
Tree contains parsed value with tag. `tree.Child(tag).V` returns value that correlated with specified tag.

Tag is defined while making parser function. There is special function `abnf.K(Rule, tag)` that make Special Rule that define tag.

## Example
This sample parse FQDN.

```
Identity ::= label *("." label)
label    ::= ALPHANUM *ldhstr
ldhstr   ::= ALPHANUM / ("-" ALPHANUM)
```
```
const (
	idFQDN int = iota
)

type Identity string

func ParseIdentity(str string) (id Identity, e error) {
	t := abnf.ParseString(str, _identity())
	if t == nil {
		e = fmt.Errorf("Invalid id text")
	} else {
		id = Identity(t.Child(idFQDN).V)
	}
	return
}

func _identity() abnf.Rule {
	return abnf.C(_fqdn(), abnf.ETX())
}

func _fqdn() abnf.Rule {
	return abnf.K(abnf.C(_label(), abnf.R0(abnf.C(abnf.V('.'), _label()))), idFQDN)
}

func _label() abnf.Rule {
	return abnf.C(abnf.ALPHANUM(), abnf.R0(_ldhstr()))
}

func _ldhstr() abnf.Rule {
	return abnf.A(abnf.ALPHANUM(), abnf.C(abnf.V('-'), abnf.ALPHANUM()))
}
```

## Notice
This parser is LL parser.
You should modify somae of ABNF definition to apply them this package.

## License
MIT

### ebnf grammar of o.t. language

*golang ebnf-parser special syntax of escape characters and so on*

The `Object` is start of grammar.

````ebnf
Object = Qid [(":" | "::") Content {WhiteSpace Content} ";"].

Qid = [Tpl "~" ] Cls ["(" Id ")"].

Content = Object | Value.

Value = Trilean | Number | Inf | String | Object | Character.

String = SubString [{":" SubString}].

SubString = (Quote {"unicode letter"} Quote) | Character.

Character = digit {hexDigit} "U".

Quote = "`" | "'" | `"`.

Inf = ["-"]"inf".

Trilean = "true" | "false" | "null".

Number = ["-"] digit {digit} (["." digit | {digit}] | {hexDigit} "H").

Tpl = ident.

Cls = ident.

Id = ident.

ident = "unicode letter" | "$" [digit | "unicode letter" | "@" | "#" | "$" | "%" | "^" | "&" | "*" | "-" | "_" | "=" | "+" | "," | "." | "?" | "!" | "/" | "|" | "\"].

digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9".

hexDigit = "A" | "B" | "C" | "D" | "E" | "F" | digit.

space = " " | "\n" | "\r" | "\t".

WhiteSpace = space {space}.

````

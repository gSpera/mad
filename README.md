MAD
===

MaD is a markdown preprocessor, it aims to create the best writing enviroment for every use case.

MaD works with little scripts that extendends the Markdown language
```
[math]: #(2x+ 2)
```
For example calls the `math` script with the argument `2x + 2`

The location of all scripts in by default `$HOME/.config/mad/bin/` but can be changed by `MAD_PATH`

The script is called with the first argument being all the given argument and then the same argument splitted by spaces.
```
$HOME/.config/mad/bin/math "2x + 2" 2x + 2
```

### Parameters

The script can access various enviroment variable to adujst it's output
Name | Type | Desc
-----|------|-----
`MAD_ISPREVIEW` | Bool | If true the script may print a textual interpretation of the input that a linter may display as a preview
`MAD_ISBLOCK` | Bool | If true the input was ia a **block** mode
`MAD_FULLINPUT` | String | The full input not splitted
`MAD_INPUTLEN` | Int | The input len

The boolean type is either "true" or "false"


### Roadmap

- [X] Preview
- [ ] Multiple Output
- [ ] Better parser

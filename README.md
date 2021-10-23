---
title: "Automated Programming Language Identification"
tags: ["Python", "Vale"]
---

# pb

Write once, publish anywhere

| Syntax      | Description | Test Text     |
| :---        |    :----:   |          ---: |
| Header      | Title       | Here's this   |
| Paragraph   | Text        | And more      |

This is the start.

$$
    \begin{matrix}
    1 & x & x^2 \\
    1 & y & y^2 \\
    1 & z & z^2 \\
    \end{matrix}
$$

Here is a p.[^1]

This is the 2nd to last node.

$$
\iiint_V \mu(u,v,w) \,du\,dv\,dw
$$

This is a CS thing $2^32$ that should be nice.

The math is $\zeta(s) = \sum 1/n^{s}$ good too.

Here is a p.

```python
def quicksort(array, begin=0, end=None):
    """A QuickSort implementation.
    """
    if end is None:
        end = len(array) - 1
    def _quicksort(array, begin, end):
        if begin >= end:
            return
        pivot = partition(array, begin, end)
        _quicksort(array, begin, pivot-1)
        _quicksort(array, pivot+1, end)
    return _quicksort(array, begin, end)
```

That's some text with a footnote.[^2]

```golang
var boldKeywords = styles.Register(chroma.MustNewStyle("bold-keywords", chroma.StyleEntries{
    chroma.Keyword:          "bold",
    chroma.KeywordConstant:  "bold",
    chroma.KeywordNamespace: "bold",
    chroma.KeywordType:      "bold",
    chroma.OperatorWord:     "bold",
}))
```

<p name="foo">HEY HO</p>

[^1]: McLaughlin, Raoul. (2010) Rome and the Distant East: Trade Routes to the ancient lands of Arabia, India, and China
[^2]: And that's the footnote.
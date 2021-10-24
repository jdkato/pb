# `pb`: Write once, publish anywhere

`pb` is a command-line tool designed to facilitate a multi-platform, scientific
publishing workflow.

It allows you to write your posts locally using scientific-styled Markdown
(tables, code blocks, footnotes, and math formulas) *and* publish to external
hosts[^1] while maintaining the structure of the original post.

## Installation

If you have the [Go toolchain][4] installed, then you can use `go get`:

```
go get github.com/jdkato/pb
```

Otherwise, download one of the [pre-built binaries][1].

## Configuration

### Medium

Uploading content to Medium requires an [Integration token][2]. To generate one, go to your [settings page][3], select "Integration tokens" from the sidebar, and create one:

<table>
    <tr>
        <td width="50%">
            <a href="https://user-images.githubusercontent.com/8785025/138613054-9487c146-29b0-4417-a33d-c9e0af6e561a.png">
                <img src="https://user-images.githubusercontent.com/8785025/138613054-9487c146-29b0-4417-a33d-c9e0af6e561a.png" width="100%">
            </a>
        </td>
        <td width="50%">
            <a href="https://user-images.githubusercontent.com/8785025/138613173-ca33e57b-9248-49c0-b472-65eeb31278eb.png">
                <img src="https://user-images.githubusercontent.com/8785025/138613173-ca33e57b-9248-49c0-b472-65eeb31278eb.png" width="100%">
            </a>
        </td>
    </tr>
    <tr>
        <td width="50%">
          Select 'Integrations tokens' from the sidebar.
        </td>
        <td width="50%">Create an Integration token; the description can be anything you want.</td>
    </tr>
</table>

Then, enter the `pb configure` command, follow the instructions, and enter your token. For math formulas on Medium, you also need to have [Inkscape][5] installed and available on your `$PATH`.

## Usage

```
pb - A multi-platform publishing workflow.

Usage:	pb [options] [command] [arguments...]
	pb --to medium file.md
	pb configure

pb is a tool for cross-posting Markdown content while preserving structural
elements (math typesetting, syntax highlighting, diagrams, etc.) across
multiple platforms.

Flags:

 -h, --help       Print this help message.
 -d, --image-dir  Search directory for local images.
 -t, --to         Comma-delimited list of destination platforms.
 -v, --version    Print the current version.

Commands:

 configure        Run an interactive configuration wizard.
 ```
 
 The basic command is
 
 ```
 pb -d <image path> <markdown file>
 ```
 
 Where `<image path>` is the directory where your local images are stored. For example, if you have an image definition like 
 
 ```markdown
 ![A demo of uploading content to Medium](/img/medium-upload.gif)
 ```
 
 and the file is stored at `/some/path/static/img/medium-upload.gif`, then you'd use:
 
  ```
 pb -d /some/path/static <markdown file>
 ```

[1]: https://github.com/jdkato/pb/releases
[2]: https://help.medium.com/hc/en-us/articles/213480228-Get-an-integration-token-for-your-writing-app
[3]: https://medium.com/me/settings
[4]: https://golang.org/
[5]: https://inkscape.org/

[^1]: Currently, only [Medium](https://medium.com/) is supported. Other hosts, like DEV and Hashnode, will be added in a future release. 

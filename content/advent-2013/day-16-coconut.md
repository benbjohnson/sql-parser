+++
title = "Go Advent Day 16 - Coconut: a pure Go blogging engine"
date = 2013-12-16T06:40:42Z
author = ["Micah Nordland"]
series = ["Advent 2013"]
+++


## Quickstart

[Coconut](https://github.com/mpnordland/coconut) is a simple blogging engine. 
It has two kinds of content: Articles and Pages.

- Articles are stored in the articles directory and can be requested by using their file name (minus the required ".md" file ext) as the first and only part of the path.
- Pages have their url and file paths specified in the config file. All page file paths are relative to the static directory, but this will probably change.

Articles have a number of meta data fields that you can use, here is an example article:

    // the meta data is inside the '---'
    ---
    title: Your First Post!
    author: You
    tags: 
     - first
     - new
    date: "Oct 24 2013 13:05"
    ---
    
    This is your first post! You can use markdown to format this.
    Markdown is a format designed for writers! coconut uses markdown for any content you want to enter.
    Articles go in the `articles` directory in your coconut installation directory.
    
    Technically, you can put pages anywhere, you just have to specify the url path for the page and the file path in conf.yaml.
    
    Here's some examples of what you can do:
        
        - *italics* by putting '\*' around some words 
        - **bold** by putting '\*\*' around words
        - lists by starting each line with a dash and a space '- '
        
        You can do code by enclosing the block with three backtick '\`\`\`'
        or by indenting each line with four spaces.
        For inline `code` you can enclose the snippet with single backtick '\`'
        
        ```
            func test() {
                        5+4
            }
            //see?
        ``` 
The meta data is then made available later in the template system (see below) with the caveat that the field name become capitalized. This is a result of the template and YAML libraries requiring exported fields on structs. Page and article bodies both use markdown for formatting but it is legal to use regular HTML as well per the 
[Markdown spec](http://daringfireball.net/projects/markdown/syntax#html).

The theme system is based on layered templates. There are templates for articles and pages along with a template containing the global layout. Templates look like this:

    <a class="header-link" href="/{{Path}}"><h2>{{Title}}</h2></a>
    <p>by {{Author}} on {{Date}}</p>
    
    {{{Body}}}
    
    <div class="tags">
    Tags: {{#Tags}}<a href="/tag/{{.}}">{{.}}</a> {{/Tags}}
    </div>
    
    {{#FullView}}
    {{! Place anything you want to only show when a user is viewing a blog post directly. }}
    {{/FullView}}

This is the article template. It has the most fields, the other templates have just one. The global template has a "Content" field which is used to hold whatever content you want to display. The page template just has a title field. 

The theme system uses the mustache templating library which in this case makes for cleaner code by specifying which fields are unescaped in the templates. Note the triple curly braces surrounding the Body field. That is the syntax for specifying escaping.

Once you've got Coconut up and running, don't hesitate to tweak these, Coconut watches these files and immediately updates them.

Coconut should be fairly cross platform, I haven't tested it on anything other than Linux, but as there is nothing explicitly *nix only, I would expect it to run fine on Windows and Mac.

If you are running it alongside IIS you may want to have Coconut served through that via FCGI by setting `protocol` to `"fcgi"` in your `conf.yaml`.

There is a system to upload articles which can be accessed at `/publish`. The user name and password are configured in conf.yaml in the users section. The definition is like this:

    admin: "$2a$12$aiTmsda0ZUcjWJ5PWqDnvOzONwRiisZUZGLv.fGF.SfT2O/9mG6Pi"

As you can see, the password is stored in bcrypt format. You can hash a password of your choice with the included `bcrypt_password_hash` program found in the directory of the same name. The work factor can be specified by the `-wf` option. I recommend a value of 12. 

Take care as the upload system will happily overwrite any articles that have the same name as the file you upload. That's pretty much all you need to know to use Coconut.

Enjoy!

## The Story

This section is for those of you who, like me, want to know the back story of a particular piece of code. 

Coconut is a sort of "kill two birds with one stone" project for me. I created it to learn the basics of how a web app works, and to get a blog for my [personal site](http://www.thoughtfuldragon.com/).

As often happens with my projects, the requirements started out rather nebulous: articles are files containing markdown, with some sort of post list on the front page. 
Unlike some of my other projects, these basic requirements actually survived through to the working program. I hadn't done anything web related in Go yet, so I went and searched for a Go web framework. 

The first one I found was [revel](http://robfig.github.io/revel/). Revel is fairly full featured, but me being a noob, I couldn't get the production profile to work. That frustration, led me to seek something simpler, which I found in [web.go](http://webgo.io/). Say what you may about it being old and such, web.go has a much better API than the other Go web frameworks that I've seen. 

It's not perfect, and lacks some stuff but I found it to fit very well to what I was trying to do. Because it uses plain functions for handling the request, it was very easy to get up and running. 

The first version was implemented entirely with functions. It was pretty messy. Once I had something I felt was interesting enough to throw up on GitHub, I went back and did some clean up. 

Splitting related functions out into separate files and grouping them with structs made the code a little cleaner. At that point I had four things that I thought needed improving/implementing: more flexible theming, pages (e.g. an about page), a list of posts on the front page, and meta data for articles. 

Some of those things were easy, but the theming improvements and the list of posts proved difficult to implement, so I rewrote the content and theming systems.
Now working with them is much easier.

Coconut runs fairly fast, fast enough for my needs anyways. I run it on a Raspberry Pi through [Pagekite](https://pagekite.net) and I get around one second load times. On my desktop over localhost pages load in ~70 microseconds. Basic analysis seems to indicate there is some fairly low hanging performance tweaks that can be made, such as gziping stuff.

I'm very pleased with what I was able to do. Go is an impressive language, with some great features. You can't write an OS with it yet, but I think that's only a matter of time.

---
Name: Bookmark Design
Version: 0.1
---
# Bookmark

## Introduction

### What is Bookmark?

Bookmark is two things:

1.  A system for creating books from folders of markdown files
2.  An open standard that other—better—tools might adopt.

Bookmark (the app) is a preprocessor for markdown and a system for keeping [Markdown][1] projects organised. It's primarily aimed at people compiling a book, but it could be used for a variety of other projects. It's not, however, intended to be for websites or blogs—bookmark aims to spit out one final piece of markdown. It is intended as both a useful tool, and a reference implementation of the Bookmark standard.

Bookmark (the open standard) is a description of a way to combine a folder of markdown files into a single formatted document. It is hoped that better, more user-friendly, tools will adopt this standard.

[1]: https://daringfireball.net/projects/markdown/ "John Gruber's Markdown"

### What is the difference between Bookmark and Markdown?

Markdown is a system for converting a plain text document into a formatted document. Bookmark is a system for converting a directory of multiple Markdown documents into a single formatted manuscript.

### Basic idea

The basic idea of Bookmark is to put a bunch of Markdown files into a folder and treat them as one continuous document, perhaps one file per chapter, and then output them as a single formatted document.

There are major issues with this idea:

1.  Folders have no inherent order. You can sort them alphabetically, or by date, but they have no way to manually order them.
2.  Folders have no way of adding metadata to them. This makes it tricky to add a title to a folder of files, unless you use the folder name.

Files have similar problems:

-   Markdown files, being plaintext files, have no standard way to store metadata either, which makes something as simple as adding an author's name to a project a little tricky.

Even if there was a simple way to order and add a title, there are a lot of other things you need to create a book—most of which involve bringing order to a large amount of information spread across a large number of files. Bookmark seeks to help with this too.

### Goals

The goals of Bookmark are:

1. Create a system to compile many markdown files into one markdown document
2. Provide tools to help in this process
3. Be a tool for writers to use—not a tool primarily for techies
4. Keep everything you need to make a document in one, portable, folder
5. Only use editable text files—no hidden files
6. Don't get in the way of established text editors
7. Don't get in the way of Git and other version control systems
8. Embrace extensibility
9. Embrace standards
10. Be open

**Non-goals:**

1. Don't try to be a blogging tool—there are enough of those already
2. Don't be a project management tool. We don't need another todo list

#### Design goals

Bookmark should respect the main tenets of Markdown:

-   Everything textual should be as readable as possible
-   Everything should be publishable as-is as plain text

So, when adding things around Markdown, we need to go for the solutions that are obvious, easy-to-read, and easy-to-work-out if you have never seen a Bookdown folder before solution.

### Similar ideas

You may have spotted that this is similar in concept to static website generators, like [Jekyll][jekyll] and [Hugo][hugo]. It is, and we have used a few of Jekyll's core ideas (like front-matter). However, Jekyll is specifically aimed at website and blog production, not books.

Bookmark also has roots in the wonderful [iA][ia] Writer and its spiritual successor [Ulysses][ulysses] (and their Granddaddy: [Scrivener][scrivener]). iA Writer and Ulysses take Markdown and create formatted text. Ulysses and Scrivener apply order to folders. Maybe, if iA had supported folder ordering, Bookmark would never have been needed.

[jekyll]: https://jekyllrb.com
[hugo]: https://gohugo.io
[ia]: https://ia.net
[ulysses]: https://ulysses.app/
[scrivener]: https://www.literatureandlatte.com/scrivener/overview

## Bookmark system

Bookmark adds 3 main things to Markdown

1. Folder metadata
2. File metadata
3. A superset of markup

What Bookmark does not supply:

- A markdown renderer.

Bookmark compiles to a single Markdown file that can then be rendered by a markdown renderer.

### File endings

Bookmark offers two styles of file extension to match the two common styles of markdown file extension. They are equivalent so pick whichever you prefer, or mix and match if you are wild and carefree.

	.bm
	.bookmark

#### Use of Markdown file extensions

Being that most text editors will not recognise bookmark files, bookmark will work with Markdown files. 

	.md
	.markdown

However, it is important to note that, once you add Bookdown-specific markup, these files will no longer be in Markdown format.

### Folder metadata

The main purpose of folder metadata is to specify the order of files in a folder. This data is held in a list in a [Sidecar file][sidecar file]. The default name for this file is 

	_index.bm

We chose this name, despite the ugly underscore, so it will always appear at the top of file listings. However you can change this filename for your project.

To add an order to a folder you list the files in the order you want as markdown list.

A folder with 4 files would alphabetically order like this:

	Chapter Four.md
	Chapter One.md
	Chapter Three.md
	Chapter Two.md

But can be ordered by adding a _folder.bm file

	- Chapter One.md
	- Chapter Two.md
	- Chapter Three.md
	- Chapter Four.md

We deliberately chose an unordered list format for two reasons:

1. It’s easy to move them around and delete entries without the numbers getting ugly
2. It’s a list

To remove a file from the listing without it reaapearing every time you automatically update the indexes, strikethough it’s listing in the _folder.bm file.

	- Chapter One.md
	- ~~Chapter Two.md~~
	- Chapter Three.md
	- Chapter Four.md

[Sidecar file]:https://en.wikipedia.org/wiki/sidecar_file

#### Folder Front Matter

Folder sidecar files can contain metadata definitions in a front-matter section. They are useful for storing values that can be used to store information about all the files in the folder, or to provide values to be used in all the files in the folder, for instance author name, latest version, dates etc.

For more about Front Matter, see below.

### File metadata

File metadata can be stored in the Front Matter portion of the file. File metadata is useful for storing metadata about the file (version, status etc) but it can also be used to define values that can be used as content within it (or other files).

For more about Front Matter, see below.

### Superset of markup

Bookmark adds a few new markup tags on top of the regular Markdown syntax.

All of Bookmark’s markup tags are denoted by square brackets. 

They look like this:

	[@variable] [#tag] [+caps]Make me bigger[-caps]

The contents of the markup tags always begin with punctuation tag to ensure they don’t class with any Markdown syntax.

They take two forms:

1. Inserts
2. Spans


#### Inserts

An *insert* is a single tag make single point in a file to insert text into the Markdown text (e.g. a date) or  to mark a place for the system to find at some point (e.g. a hashtag).

Using an insert looks like this:

	This document was generated at: [@date], [@time]

Inserts are, mostly, used for inserting the value of a variable Ito the markdown. These can be defined by the system (e.g. a date or a filename) or in the *metadata cascade* of folder metadata and Front Matter  from the Project file, each parent folder and the Front Matter from the file itself. Inserts can also be used to insert calculations or the results of dynamic searches.

Built-in values to insert include:

- date: when the output document was generated
- lastmod: when this file was last edited
- created: when this file was created
- filename: the name of the file being processed

#### Spans

*Spans* mark a section of text with a beginning and end. They are used to transform a section of text or to mark a span of text that can be queried for.

A span looks like this:

	[+scaps]This text should be transformed into small-caps[-scaps]
	
A span can	 also take a list of arguments, e.g.

	[+tag draft revisit]Text needing a revision…[-tag]

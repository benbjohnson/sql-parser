+++
title = "Bleve: Text Search Powered by Go"
date = 2014-11-16T00:00:00Z
author = ["Marty Schoch"]
tags = ["bleve","search","text"]
series = ["Birthday Bash 2014"]
+++

[Bleve](http://www.blevesearch.com/), a text indexing and search library for Go, is excited to be part of the celebration of the 5th birthday of Go.  Building a complete solution for indexing and search text requires integrating many different technologies.  Go, the language, the standard library, the whole ecosystem, combine to make this not only possible, but fun.
<!--more-->

Let's start by taking a look at a quick indexing example in Bleve:

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("myExampleIndex.bleve", mapping)

	msg := map[string]string{
	    "body": "Happy Birthday Go!",
	}
	err = index.Index("birthday", msg)

## Standard Library Text Support

The first thing Bleve does when indexing the document above is perform a series of transformations on the text in the document.  Bleve focuses on Unicode text in a UTF-8 encoding, so having clean support for this in the standard library is essential.  The `unicode/utf8` package allows us to interpret strings and byte slices as sequences of Unicode code points (runes in Go).  Then the `unicode` package provides lower-level support for working with the properties of individual runes.  These packages allow us to implement basic filters on the text, like lower-casing all of the text, or filtering words with too many or too few runes.


Often we need to use some of the more advanced features of Unicode like normalization and segmentation.  Unicode strings often have multiple representations that are semantically equivalent.  For example, consider the strings `café` and `cafe\u0301`. In the second form `\u0301` is combining character which modifies the preceding charater with the acute accent.  Normalization is a set of processes that allow us to convert and compare these representations.  Normalization is available in a separate [go.text](http://godoc.org/code.google.com/p/go.text) repository.  

Segmentation is the process of splitting text at word and sentence boundaries.  Not all languages separate words with whitespace, and punctuation within sentences and words makes this non-trivial.  Bleve has built its own Unicode [segmentation](https://github.com/blevesearch/segment) library.  But the key is that right foundations are in place.  When we did have to build custom functionality, we were able to build on top of data-structures and functions of the standard library.

## Integration with C/C++

Another often overlooked feature in Go is the ability to easily integrate with existing C/C++ libraries.  Now at Bleve, our goal is to deliver all the core functionality with a pure Go implementation.  But as a new project with a long list of features that users expect to work out of the box, we were able to bootstrap our feature list by adding tactical integrations with existing C libraries.  Integration with [libstemmer](http://snowball.tartarus.org/) gave Bleve users the ability to correctly process text in a large number of languages long before we could deliver all the functionality ourselves.  

<table width="100%">
	<tr>
		<td>
			<ul>
				<li>Danish</li>
				<li>Dutch</li>
				<li>English</li>
				<li>Finnish</li>
				<li>French</li>
			</ul>
		</td>
		<td>
			<ul>
				<li>German</li>
				<li>Hungarian</li>
				<li>Italian</li>
				<li>Norwegian</li>
				<li>Portuguese</li>
			</ul>
		</td>
		<td>
			<ul>
				<li>Romanian</li>
				<li>Russian</li>
				<li>Spanish</li>
				<li>Swedish</li>
				<li>Turkish</li>
			</ul>
		</td>
	</tr>
</table>

Another critical integration with C/C++ is in an implementation of Bleve's pluggable storage layer. The first implementation offered was an integration with [LevelDB](https://github.com/google/leveldb).  Building on top of a proven and fast key/value store let us save time and focus on search features.

*NOTE*: Subsequently, we've also added support for [Bolt](https://github.com/boltdb/bolt) a great pure Go alternative.

## Concurrency

We can't talk about Go without mentioning concurrency.  Indexing and searching text efficiently requires balancing a lot of activities at the same time.  First, we’re writing bits to disk, so we want to keep the disk as busy as we can.  Second, we’re doing (sometimes intensive) analysis of the text, so we want to keep all the cores busy as well.  On top of this, users want to be able to tune the behavior, grouping documents into larger batches for better throughput, and smaller batches for lower latency.

![](/postimages/bleve-text-search-powered-by-go/bleve-concurrency.png)

Bleve allows users of the library to configure a pool of analysis workers.  Text analysis on documents can be done in parallel as each document is independent.  Index updates within the same batch are then merged together and applied.  Currently it is up to the application to build batches of the desired size, but in the future we plan to support automatic batches to balance latency and throughput within the library.

Intelligent use of goroutines in established patterns make it possible to deliver this while keeping the code in a clean maintainable form.

## Conclusion

Together, the language, the standard library, and the exploding ecosystem around Go make Bleve a project that at its core is fun to work on.  To learn more, join [the community](https://groups.google.com/forum/#!forum/bleve) and help improve [Bleve](http://www.blevesearch.com/).
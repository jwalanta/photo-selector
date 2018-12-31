# Photo Selector

Fast, cross-platform, web-based photo selector.

## Background

I need to select photos in a folder, *fast* (note the emphasis). There are bunch of applications that can do this, like Lightroom. But I want to do it without having to import the photos, or well, pay for the application. Just a quick and simple photo selector. And it has to be cross-platform. Did I mention that it has to be *fast*?

I'm pretty comfortable with opening up a shell window, so here's the idea:

* Open shell, go to the folder containing the photos
* run the photo selector command
    * the command will start a web server
    * the photos will be resized in the background so it'll load quick
* Go to browser, open the link given by the command line
* Browse through the photos, select photos you like (keyboard shortcuts for everything)
* Copy the list of selected filenames
* Do whatever you want with that list
* Profit?? :)

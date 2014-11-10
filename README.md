gopheracademy-web
=================

Public Website for GopherAcademy.com

# Working with gopheracademy-web

1. Clone the repo
2. Install [Hugo](http://hugo.spf13.com)

## Modifying the Layout

Theme/HTML are in `layouts/` and static assets(CSS/JS) are in `static/`. When Hugo runs, the final layout is generated and served from the `public/` folder. In order to benefit from LiveReload while making changes to the site run Hugo as:

	hugo --watch server

## Viewing the site locally

In the gopheracademy-web directory:

    hugo server


View the url that Hugo provides in a browser

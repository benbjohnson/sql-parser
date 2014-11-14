gopheracademy-web
=================

There are two parts to the code now. The main Hugo app still runs the blog with **config.toml** while an aditional **config-main.toml** runs the main website.

Stay tuned while I get the layouts and design done for the main website. Right now it looks exactly like the blog page.

# Working with Hugo

1. Clone the repo
2. Install [Hugo](http://hugo.spf13.com)

## Working with the blog 

Theme/HTML are in `layouts/` and static assets(CSS/JS) are in `static/`. When Hugo runs, the final layout is generated and served from the `public/` folder. Run the Hugo server with:

	hugo --watch server --config="config.toml"

Remember to include the config flag as the main site uses another config file.

## Working with the main site

The website runs on the same Hugo app but has a few things configured differently. Layouts are in `layouts-main/` and generated pages are in `public-main/`. It uses the same `content/` folder from the blog app to pull info to the main site. Please see `config-main.toml` to understand how it's configured. 

To run the server, include the mentioned config file as a flag:

In the gopheracademy-web directory:

    hugo --watch server --config="config-main.toml"


To view the site, visit the link provided by Hugo, usually `http://localhost:1313`.

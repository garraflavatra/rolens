const esbuild = require('esbuild');
const fs = require('fs');
const path = require('path');
const sass = require('sass');

const navigationPlugin = require('@11ty/eleventy-navigation');
const { EleventyHtmlBasePlugin: basePlugin } = require("@11ty/eleventy");
const { glob, globSync } = require('glob');

const markdown = require('markdown-it')({
  html: true,
  breaks: true,
  linkify: true,
  typographer: true,
});

// Add target="_blank" to external links
const defaultAnchorRenderer = markdown.renderer.rules.link_open || function(tokens, idx, options, env, self) {
  return self.renderToken(tokens, idx, options);
};
markdown.renderer.rules.link_open = function(tokens, idx, options, env, self) {
  const url = tokens[idx].attrGet('href');
  if (url?.startsWith('http')) {
    const aIndex = tokens[idx].attrIndex('target');
    if (aIndex < 0) {
      tokens[idx].attrPush(['target', '_blank']);
    }
    else {
      tokens[idx].attrs[aIndex][1] = '_blank';
    }
    }

  return defaultAnchorRenderer(tokens, idx, options, env, self);
};

const indir = path.join(__dirname, '/../docs');
const outdir = path.join(__dirname, '/../website-dist');

function buildStyles() {
   const { css } = sass.compile(__dirname+'/sass/styles.scss', {
    loadPaths: [__dirname+'/sass'],
    style: 'compressed',
  });

  fs.mkdirSync(outdir, { recursive: true });
  fs.writeFileSync(outdir+'/styles.css', css);
}

function buildScript() {
  esbuild.buildSync({
    entryPoints: [__dirname+'/js/script.js'],
    outdir: outdir,
    minify: true,
    target: 'es2016',
  });
}

function copyStatic() {
  fs.mkdirSync(outdir, { recursive: true });
  fs.readdirSync(__dirname+'/static').forEach(function(fname) {
    fs.copyFileSync(__dirname+'/static/'+fname, outdir+'/'+fname);
  });
  globSync(indir+'/**/*.{png,jpg,jpeg}').forEach(function(fname) {
    const dest = outdir+'/'+fname.slice(indir.length + 1);
    fs.mkdirSync(path.dirname(dest), { recursive: true });
    fs.copyFileSync(fname, dest);
  });
}

module.exports = function (eleventyConfig) {
  eleventyConfig.setLibrary('md', markdown);
  eleventyConfig.addPlugin(navigationPlugin);
  eleventyConfig.addPlugin(basePlugin);
  eleventyConfig.addGlobalData('layout', 'page.liquid');

  // Static files
  eleventyConfig.on('eleventy.before', copyStatic);

  // Sass rendering
  eleventyConfig.addWatchTarget(__dirname+'/sass/');
  eleventyConfig.on('eleventy.before', buildStyles);

  // JS compilation
  eleventyConfig.addWatchTarget(__dirname+'/js/');
  eleventyConfig.on('eleventy.before', buildScript);

  // Markdown filter
  eleventyConfig.addFilter('markdowninline', function (text) {
    return markdown.renderInline(text);
  });

  // Retrieve content of a file
  eleventyConfig.addShortcode('filecontent', function (fname, a, b) {
    let str = fs.readFileSync(path.join(indir, fname)).toString();

    if (typeof a === 'number') {
      // Line indexes
      str = str.split('\n').slice(a, b).join('\n');
    }
    else if (typeof a === 'string') {
      // Anchors
      const aIndex = str.indexOf(a) + a.length;
      const bIndex = b ? str.indexOf(b, aIndex) : undefined;
      str = str.slice(aIndex, bIndex);
    }

    return str.trim();
  });

  // Global options
  return {
    dir: {
      input: indir,
      includes: '../website/includes',
      layouts: '../website/layouts',
      data: '../website/data',
      output: outdir,
    },
    templateFormats: ['html', 'liquid', 'md', '11ty.js'],
    markdownTemplateEngine: 'liquid',
    htmlTemplateEngine: 'liquid',
  };
};

module.exports = {
  eleventyNavigation: {
    key: data => data.title,
    parent: data => data.parent,
    order: data => data.order,
  },

  navigationOptions: {
    activeKey: data => data.eleventyNavigation.key,
  },

  github: {
    reportDocsIssue: data => `${data.meta.repoUrl}/issues/new?assignees=garraflavatra&labels=documentation&projects=&template=docs.yml&source=${encodeURIComponent(data.meta.siteUrl + data.page.url)}`,
    licenseUrl: data => `${data.meta.repoUrl}/blob/main/LICENSE`,
    pageSourceUrl: data => `${data.meta.repoUrl}/blob/main/${data.page.inputPath}`,
  },
};

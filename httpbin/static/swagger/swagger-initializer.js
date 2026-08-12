window.onload = function () {
  window.ui = SwaggerUIBundle({
    url: "./doc.json",
    dom_id: "#swagger-ui",
    deepLinking: true,
    validatorUrl: null,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
};

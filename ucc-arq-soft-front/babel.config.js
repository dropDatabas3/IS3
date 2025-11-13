module.exports = {
  presets: [
    [
      '@babel/preset-env',
      {
        // Target modern browsers and current Node; let Next handle polyfills
        targets: { browsers: ['defaults'], node: 'current' },
        // IMPORTANTE: no forzar modules: 'commonjs' para no romper el build de Next
        // Si Jest necesitara CommonJS, se puede configurar solo en jest.config.
      }
    ],
    [
      '@babel/preset-react',
      {
        // Use the new JSX transform (no need to import React)
        runtime: 'automatic'
      }
    ],
    [
      '@babel/preset-typescript',
      {
        allowDeclareFields: true
      }
    ]
  ]
};
module.exports = {
  presets: [
    [
      '@babel/preset-env',
      {
        // Target modern browsers and current Node; let Next handle polyfills
        targets: { browsers: ['defaults'], node: 'current' },
        // For Jest, compile ESM to CommonJS so require() works for transformed deps
        modules: 'commonjs'
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
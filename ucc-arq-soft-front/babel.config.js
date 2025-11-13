module.exports = {
  presets: [
    [
      '@babel/preset-env',
      {
        // Target modern browsers and current Node; let Next handle polyfills
        targets: { browsers: ['defaults'], node: 'current' },
        modules: false
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
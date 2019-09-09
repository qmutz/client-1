module.exports = {
  env: {es6: true},
  parser: 'babel-eslint',
  extends: ['eslint:recommended'],
  overrides: [
    {
      files: ['*.tsx', '*.ts', '*.d.ts'],
      rules: {
        'no-undef': 'off', // ts itself will catch this
        'no-unused-vars': 'off', // ts itself will catch this
      },
    },
  ],
  parserOptions: {
    ecmaVersion: 6,
    sourceType: 'module',
    ecmaFeatures: {jsx: true},
    babelOptions: {
      configFile: __dirname + '/babel.config.js',
    },
  },
  globals: {
    KB: 'readonly',
    requestAnimationFrame: 'readonly',
    cancelAnimationFrame: 'readonly',
    require: 'readonly',
    __DEV__: false,
    __STORYBOOK__: false,
    __STORYSHOT__: false,
  },
  plugins: ['promise', 'react', 'filenames', 'import', 'react-hooks'],
  settings: {
    'import/core-modules': ['electron', 'react-native'],
    'import/extensions': ['.js', '.tsx', '.d.ts', '.native.tsx', '.desktop.tsx', '.native.js', '.desktop.js'],
    'import/resolver': {
      node: {
        extensions: ['.js', '.tsx', '.d.ts', '.native.tsx', '.desktop.tsx', '.native.js', '.desktop.js'],
      },
    },
    react: {
      version: 'detect',
    },
  },
  rules: {
    'array-callback-return': 'error',
    'comma-dangle': ['error', 'always-multiline'],
    'dot-notation': 'off',
    'func-call-spacing': 'off',
    'generator-star-spacing': 'off',
    'jsx-quotes': 'off',
    'lines-between-class-members': 'off',
    'no-duplicate-imports': 'off',
    'no-empty': 'off',
    'no-extra-semi': 'off',
    'no-implied-eval': 'error',
    'no-loop-func': 'off',
    'no-mixed-operators': 'off',
    'no-script-url': 'error',
    'no-self-compare': 'error',
    'no-sequences': 'error',
    'no-shadow': 'warn',
    'no-unused-expressions': 'off',
    'no-use-before-define': 'off',
    'no-useless-return': 'off',
    'object-curly-even-spacing': 'off',
    'object-curly-spacing': 'off',
    'prefer-const': 'warn',
    'quote-props': 'off',
    'react-hooks/exhaustive-deps': 'warn',
    'react-hooks/rules-of-hooks': 'error',
    'react/boolean-prop-naming': 'error',
    'react/button-has-type': 'off',
    'react/default-props-match-prop-types': 'error',
    'react/destructuring-assignment': 'off',
    'react/display-name': 'off',
    'react/forbid-component-props': 'off',
    'react/forbid-dom-props': 'off',
    'react/forbid-elements': 'off',
    'react/forbid-foreign-prop-types': 'off',
    'react/forbid-prop-types': 'off',
    'react/jsx-boolean-value': ['error', 'always'],
    'react/jsx-child-element-spacing': 'off',
    'react/jsx-closing-bracket-location': 'off',
    'react/jsx-closing-tag-location': 'off',
    'react/jsx-curly-brace-presence': 'off',
    'react/jsx-curly-newline': 'off',
    'react/jsx-curly-spacing': 'off',
    'react/jsx-equals-spacing': 'off',
    'react/jsx-filename-extension': 'off',
    'react/jsx-first-prop-new-line': 'off',
    'react/jsx-fragments': ['error', 'syntax'],
    'react/jsx-handler-names': 'off',
    'react/jsx-indent': 'off',
    'react/jsx-indent-props': 'off',
    'react/jsx-key': 'error',
    'react/jsx-max-depth': 'off',
    'react/jsx-max-props-per-line': 'off',
    'react/jsx-no-bind': 'off',
    'react/jsx-no-comment-textnodes': 'error',
    'react/jsx-no-duplicate-props': 'error',
    'react/jsx-no-literals': 'off',
    'react/jsx-no-target-blank': 'error',
    'react/jsx-no-undef': 'error',
    'react/jsx-one-expression-per-line': 'off',
    'react/jsx-pascal-case': 'off',
    'react/jsx-props-no-multi-spaces': 'off',
    'react/jsx-props-no-spreading': 'off',
    'react/jsx-sort-default-props': 'off',
    'react/jsx-sort-props': 'off',
    'react/jsx-space-before-closing': 'off',
    'react/jsx-tag-spacing': 'off',
    'react/jsx-uses-react': 'error',
    'react/jsx-uses-vars': 'error',
    'react/jsx-wrap-multilines': 'off',
    'react/no-access-state-in-setstate': 'error',
    'react/no-array-index-key': 'off',
    'react/no-children-prop': 'off',
    'react/no-danger': 'error',
    'react/no-danger-with-children': 'error',
    'react/no-deprecated': 'warn',
    'react/no-did-mount-set-state': 'warn',
    'react/no-did-update-set-state': 'warn',
    'react/no-direct-mutation-state': 'error',
    'react/no-find-dom-node': 'warn',
    'react/no-is-mounted': 'error',
    'react/no-multi-comp': 'off',
    'react/no-redundant-should-component-update': 'error',
    'react/no-render-return-value': 'error',
    'react/no-set-state': 'off',
    'react/no-string-refs': 'error',
    'react/no-this-in-sfc': 'error',
    'react/no-typos': 'warn',
    'react/no-unescaped-entities': 'off',
    'react/no-unknown-property': 'error',
    'react/no-unsafe': 'error',
    'react/no-unused-prop-types': 'error',
    'react/no-unused-state': 'error',
    'react/no-will-update-set-state': 'warn',
    'react/prefer-es6-class': 'off',
    'react/prefer-read-only-props': 'error',
    'react/prefer-stateless-function': 'warn',
    'react/prop-types': 'off',
    'react/react-in-jsx-scope': 'error',
    'react/require-default-props': 'error',
    'react/require-optimization': 'off',
    'react/require-render-return': 'off',
    'react/self-closing-comp': 'off',
    'react/sort-comp': 'off',
    'react/sort-prop-types': 'off',
    'react/state-in-constructor': 'off',
    'react/static-property-placement': 'off',
    'react/style-prop-object': 'error',
    'react/void-dom-elements-no-children': 'error',
    'sort-keys': ['error', 'asc', {caseSensitive: true, natural: false}],
    'yield-star-spacing': 'off',
    camelcase: 'off',
    curly: 'off',
    indent: 'off',
    quotes: 'off',
    strict: ['error', 'global'],
  },
}

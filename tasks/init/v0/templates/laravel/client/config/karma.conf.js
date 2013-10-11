
module.exports = function(config) {
  config.set({
    // base path, that will be used to resolve files and exclude
    basePath: '..',

    frameworks: ['jasmine'],

    // list of files / patterns to load in the browser
    files: [
      'app/components/jquery-1.8.3/jquery.js',
      'app/components/bootstrap/js/transition.js',
      'app/components/bootstrap/js/modal.js',
      'app/components/bower-angular/angular.js',
      'app/components/bower-angular/angular-*.js',
      'app/components/angular-bootstrap/ui-bootstrap-tpls.js',
      'app/scripts/*.js',
      'app/scripts/**/*.js',
      'app/scripts/mocks/*.js',
      'test/unit/**/*.js'
    ],

    // list of files to exclude
    exclude: ['app/components/bower-angular/angular-scenario.js', 'app/scripts/vendor/**/*.js'],

    // use dots reporter, as travis terminal does not support escaping sequences
    // possible values: 'dots', 'progress', 'junit', 'teamcity'
    // CLI --reporters progress
    reporters: ['progress'],

    // enable / disable colors in the output (reporters and logs)
    // CLI --colors --no-colors
    colors: true,

    // enable / disable watching file and executing tests whenever any file changes
    // CLI --auto-watch --no-auto-watch
    autoWatch: true,

    // Start these browsers, currently available:
    // - Chrome
    // - ChromeCanary
    // - Firefox
    // - Opera
    // - Safari (only Mac)
    // - PhantomJS
    // - IE (only Windows)
    // CLI --browsers Chrome,Firefox,Safari
    browsers: [],

    // If browser does not capture in given timeout [ms], kill it
    // CLI --capture-timeout 5000
    captureTimeout: 5000,

    // Auto run tests on start (when browsers are captured) and exit
    // CLI --single-run --no-single-run
    singleRun: false,

    // report which specs are slower than 500ms
    // CLI --report-slower-than 500
    reportSlowerThan: 500,

    // compile coffee scripts
    preprocessors: {
    },

    plugins: [
      'karma-jasmine',
      'karma-phantomjs-launcher'
    ]
  });
};


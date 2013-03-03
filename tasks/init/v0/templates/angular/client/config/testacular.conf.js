
basePath = '..';

files = [
  JASMINE,
  JASMINE_ADAPTER,
  'app/components/jquery/jquery.js',
  'app/components/bootstrap/js/bootstrap-transition.js',
  'app/components/bootstrap/js/bootstrap-modal.js',
  'app/components/bower-angular/angular.js',
  'app/components/bower-angular/angular-*.js',
  'app/scripts/*.js',
  'app/scripts/**/*.js',
  'test/unit/**/*.js'
];

exclude = ['app/components/bower-angular/angular-scenario.js']

autoWatch = true;

browsers = [];

junitReporter = {
  outputFile: 'test_out/unit.xml',
  suite: 'unit'
};

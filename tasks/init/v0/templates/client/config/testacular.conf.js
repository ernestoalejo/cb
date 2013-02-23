
basePath = '..';

files = [
  JASMINE,
  JASMINE_ADAPTER,
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

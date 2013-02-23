
basePath = '..';

files = [
  JASMINE,
  JASMINE_ADAPTER,
  'dist/scripts/*.js',
  'app/components/bower-angular/angular-mocks.js',
  'test/unit/**/*.js'
];

autoWatch = true;

browsers = [];

junitReporter = {
  outputFile: 'test_out/unit.xml',
  suite: 'unit'
};

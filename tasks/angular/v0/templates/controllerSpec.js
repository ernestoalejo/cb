{{ if not .Exists }}'use strict';

{{ else }}
{{ end }}
describe('Controller: {{ .Data.Name }}', function() {
  beforeEach(module('{{ .Data.Module }}'));

  var scope;
  beforeEach(inject(function($injector) {
    var $controller = $injector.get('$controller');
    var $rootScope = $injector.get('$rootScope');

    scope = $rootScope.$new();
    $controller('{{ .Data.Name }}', {$scope: scope});
  }));
});

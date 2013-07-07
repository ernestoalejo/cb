{{ if not .Exists }}'use strict';

{{ else }}
{{ end }}
describe('Controller: {{ .Data.Name }}', function() {
  beforeEach(module('controllers.{{ .Data.Module }}'));

  var scope;
  beforeEach(inject(function($controller, $rootScope) {
    scope = $rootScope.$new();
    $controller('{{ .Data.Name }}', {$scope: scope});
  }));
});

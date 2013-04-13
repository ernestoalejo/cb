{{ if not .Exists }}'use strict';

{{ else }}
{{ end }}
describe('Service: {{ .Data.Name }}', function() {
  beforeEach(module('services.{{ .Data.Module }}'));

  var {{ .Data.Name }};
  beforeEach(inject(function($injector) {
    {{ .Data.Name }} = $injector.get('{{ .Data.Name }}');
  }));
});

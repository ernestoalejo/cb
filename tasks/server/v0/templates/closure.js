
window.CLOSURE_NO_DEPS = true;
window.CLOSURE_BASE_PATH = 'http://localhost:9810/input/';

{{ . }}

for(var i = 0; i < cb_deps.length; i++) {
  goog.require(cb_deps[i]);
}

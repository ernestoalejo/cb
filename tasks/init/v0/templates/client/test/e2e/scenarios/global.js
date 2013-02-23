'use strict';


describe('Bogus Test', function() {
  beforeEach(function() {
    browser().navigateTo('/');
  });

  it('should render view2 when user navigates to /view2', function() {
    expect(true).toBe(true);
  });
});

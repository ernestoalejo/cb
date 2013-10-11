<?php namespace Admin;

use Auth;
use User;

class HomeControllerTest extends \TestCase {

  protected $admin = true;

  public function testIndex() {
    $user = new User();
    $user->id = 'foo';
    Auth::shouldReceive('user')->once()->andReturn($user);

    $this->call('GET', '');
  }

}

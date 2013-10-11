<?php

use Symfony\Component\Console\Input\ArrayInput;
use Symfony\Component\Console\Output\NullOutput;

class TestCase extends Illuminate\Foundation\Testing\TestCase {

	protected $nestedViewsData = array();
	protected $admin = false;
  protected $ignoreEvent;

  private $eventCalls = array();

	/**
	 * Creates the application.
	 *
	 * @return Symfony\Component\HttpKernel\HttpKernelInterface
	 */
	public function createApplication() {
		$unitTesting = true;
		$testEnvironment = 'testing';
		return require __DIR__.'/../../bootstrap/start.php';
	}

	public function setUp() {
		parent::setUp();

    $this->ignoreEvent = new stdClass();
		$this->setAdmin();
	}

	public function setAdmin() {		
		$host = ($this->admin) ? Config::get('site.admin-domain') : Config::get('site.front-domain');
		$this->client->setServerParameters(array('HTTP_HOST' => $host));
	}

	public function registerNestedView($view) {
	  View::composer($view, function($view){
	    $this->nestedViewsData[$view->getName()] = $view->getData();
	  });
	}

	/**
	 * Assert that the given view has a given piece of bound data.
	 *
	 * @param  string|array  $key
	 * @param  mixed  $value
	 * @return void
	 */
	public function assertNestedViewHas($view, $key, $value = null) {
		if (is_array($key)) {
			return $this->assertNestedViewHasAll($view, $key);
		}
		if (!isset($this->nestedViewsData[$view])) {
			return $this->assertTrue(false, 'The view was not called.');
		}

		$data = $this->nestedViewsData[$view];

		if (is_null($value)) {
			$this->assertArrayHasKey($key, $data);
		} else {
			if(isset($data[$key])) {
				$this->assertSame($value, $data[$key]);
			} else {
				return $this->assertTrue(false, 'The View has no bound data with this key.');            
			}
		}
	}

	/**
	 * Assert that the view has a given list of bound data.
	 *
	 * @param  array  $bindings
	 * @return void
	 */
	public function assertNestedViewHasAll($view, array $bindings) {
		foreach ($bindings as $key => $value) {
			if (is_int($key)) {
				$this->assertNestedViewHas($view, $value);
			} else {
				$this->assertNestedViewHas($view, $key, $value);
			}
		}
	}

	public function assertNestedView($view) {
	  $this->assertArrayHasKey($view, $this->nestedViewsData);
	}

  public function get($url, $data = array()) {
    $resp = $this->call('GET', $url, $data);
    return $this->decodeJSONResponse($resp);
  }

  public function post($url, $data = array()) {
    $resp = $this->call('POST', $url, array(), array(), array(), json_encode($data));
    return $this->decodeJSONResponse($resp);
  }

	public function postUnsafe($url, $data = array()) {
    $resp = $this->call('POST', $url, $data);
		return $this->decodeUnsafeJSONResponse($resp);
	}

  public function decodeJSONResponse($resp) {
    return json_decode(substr($resp->getContent(), 6), true);
  }

	public function decodeUnsafeJSONResponse($resp) {
		return json_decode($resp->getContent(), true);
	}

	public function migrate($seed = false) {
		$artisan = Illuminate\Console\Application::start($this->app);

		$status = $artisan->run(new ArrayInput(array('migrate')), new NullOutput());
		$this->assertSame(0, $status);
		if ($seed) {
			$status = $artisan->run(new ArrayInput(array('db:seed')), new NullOutput());
			$this->assertSame(0, $status);
		}
	}

  protected function callRestrictedMethod($obj, $name, array $args) {
    $class = new \ReflectionClass($obj);
    $method = $class->getMethod($name);
    $method->setAccessible(true);
    return $method->invokeArgs($obj, $args);
  }

  protected function watchEvent($event, Callable $check = null) {
    $this->eventCalls[$event] = 0;
    \Event::listen($event, function($data) use ($event, $check) {
      if (!is_null($check)) {
        $ret = $check($data);
        if ($ret === $this->ignoreEvent) {
          return;
        }
      }
      $this->eventCalls[$event]++;
    });
  }

  protected function getEventCalls($name) {
    return $this->eventCalls[$name];
  }

  protected function assertEventCalls($event, $calls) {
    if (!isset($this->eventCalls[$event])) {
      throw new \Exception("event name not watched: " . $event);
    }
    $this->assertSame($calls, $this->eventCalls[$event]);
  }

}

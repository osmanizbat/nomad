import { currentURL } from '@ember/test-helpers';
import { run } from '@ember/runloop';
import { module, skip, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';
import ServerMonitor from 'nomad-ui/tests/pages/servers/monitor';

let agent;
let managementToken;
let clientToken;

module('Acceptance | server monitor', function(hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function() {
    agent = server.create('agent');

    managementToken = server.create('token');
    clientToken = server.create('token');

    window.localStorage.nomadTokenSecret = managementToken.secretId;

    run.later(run, run.cancelTimers, 500);
  });

  test('/servers/:id/monitor should have a breadcrumb trail linking back to servers', async function(assert) {
    await ServerMonitor.visit({ name: agent.name });

    assert.equal(ServerMonitor.breadcrumbFor('servers.index').text, 'Servers');
    assert.equal(ServerMonitor.breadcrumbFor('servers.server').text, agent.name);

    await ServerMonitor.breadcrumbFor('servers.index').visit();
    assert.equal(currentURL(), '/servers');
  });

  skip('the monitor page immediately streams agent monitor output at the info level', async function(assert) {
    await ServerMonitor.visit({ name: agent.name });

    const logRequest = server.pretender.handledRequests.find(req =>
      req.url.startsWith('/v1/agent/monitor')
    );
    assert.ok(ServerMonitor.logsArePresent);
    assert.ok(logRequest);
    assert.ok(logRequest.url.includes('log_level=info'));
  });

  test('switching the log level persists the new log level as a query param', async function(assert) {
    await ServerMonitor.visit({ name: agent.name });
    await ServerMonitor.selectLogLevel('Debug');
    assert.equal(currentURL(), `/servers/${agent.name}/monitor?level=debug`);
  });

  test('when the current access token does not include the agent:read rule, a descriptive error message is shown', async function(assert) {
    window.localStorage.nomadTokenSecret = clientToken.secretId;

    await ServerMonitor.visit({ name: agent.name });
    assert.notOk(ServerMonitor.logsArePresent);
    assert.ok(ServerMonitor.error.isShown);
    assert.equal(ServerMonitor.error.title, 'Not Authorized');
    assert.ok(ServerMonitor.error.message.includes('agent:read'));

    await ServerMonitor.error.seekHelp();
    assert.equal(currentURL(), '/settings/tokens');
  });
});

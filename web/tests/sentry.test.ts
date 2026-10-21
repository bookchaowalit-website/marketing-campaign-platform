import assert from "node:assert/strict";
import { describe, it } from "node:test";
import { scrubEvent, sentryOptions } from "../src/lib/sentry.ts";

describe("sentryOptions", () => {
  it("stays disabled when no DSN is configured", () => {
    const options = sentryOptions(undefined);
    assert.equal(options.enabled, false);
  });

  it("reports errors only: no PII, tracing, trace headers or breadcrumbs", () => {
    const options = sentryOptions("https://key@o0.ingest.sentry.io/1");
    assert.equal(options.enabled, true);
    assert.equal(options.sendDefaultPii, false);
    assert.equal(options.tracesSampleRate, undefined);
    assert.deepEqual(options.tracePropagationTargets, []);
    assert.equal(options.maxBreadcrumbs, 0);
    assert.equal(options.beforeSendTransaction(), null);
  });
});

describe("scrubEvent", () => {
  it("drops request data, user, breadcrumbs and extra context", () => {
    const event = scrubEvent({
      message: "boom",
      request: { url: "https://marketing.example/campaigns?token=secret", headers: { cookie: "a=b" } },
      user: { ip_address: "1.2.3.4", email: "x@example.com" },
      breadcrumbs: [{ message: "clicked" }],
      extra: { body: "secret" },
    });
    assert.equal(event.message, "boom");
    assert.equal(event.request, undefined);
    assert.equal(event.user, undefined);
    assert.equal(event.breadcrumbs, undefined);
    assert.equal(event.extra, undefined);
  });

  it("removes local variables from stack frames", () => {
    const event = scrubEvent({
      exception: {
        values: [{ type: "Error", stacktrace: { frames: [{ function: "f", vars: { password: "p" } }] } }],
      },
    });
    const frame = event.exception?.values?.[0]?.stacktrace?.frames?.[0];
    assert.equal(frame?.function, "f");
    assert.equal(frame?.vars, undefined);
  });
});

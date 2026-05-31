"use client";

import { useMemo, useState } from "react";
import { EMAIL_TEMPLATES, type EmailTemplate } from "@/lib/email-templates";
import { injectMockEmail } from "@/lib/api";
import { humanize } from "@/lib/humanize";

/**
 * Mock-mail composer for the Lab page. Pick a template, optionally tweak,
 * send → email lands in the global inbox so process_inbox / validate_invoice
 * / route_rfi scenarios pick it up on their next run.
 *
 * Notes for the reader:
 *  - We track edits as a partial overlay over the chosen template, so
 *    switching template is "go back to defaults" cleanly.
 *  - Sends generate a stable, readable id like "mock_vendor_impersonation_<ts>"
 *    so the inbox stays scannable and re-runs don't collide.
 */
export function MockMailComposer({ onSent }: { onSent?: () => void }) {
  const [tplId, setTplId] = useState<string>(EMAIL_TEMPLATES[0].id);
  const [overlay, setOverlay] = useState<Partial<EmailTemplate["payload"]>>({});
  const [busy, setBusy] = useState(false);
  const [flash, setFlash] = useState<string | null>(null);

  const tpl = useMemo(
    () => EMAIL_TEMPLATES.find((t) => t.id === tplId) ?? EMAIL_TEMPLATES[0],
    [tplId]
  );
  const fields = useMemo(
    () => ({ ...tpl.payload, ...overlay }),
    [tpl, overlay]
  );

  const update = <K extends keyof EmailTemplate["payload"]>(
    k: K,
    v: EmailTemplate["payload"][K]
  ) => setOverlay((o) => ({ ...o, [k]: v }));

  const flashFor = (msg: string, ms = 4000) => {
    setFlash(msg);
    setTimeout(() => setFlash(null), ms);
  };

  const onPickTemplate = (id: string) => {
    setTplId(id);
    setOverlay({});
  };

  const onSend = async () => {
    if (busy) return;
    setBusy(true);
    const id = `mock_${tpl.id}_${Date.now()}`;
    const res = await injectMockEmail({ id, ...fields });
    setBusy(false);
    if (res.ok) {
      flashFor(`Sent "${fields.subject}" → inbox (${id}).`);
      onSent?.();
    } else {
      flashFor(`Send failed: ${res.error || "unknown error"}`);
    }
  };

  return (
    <section className="lab-card mock-mail">
      <h3>Mock mail</h3>
      <p className="lab-helper">
        Inject a realistic inbound email into the inbox. Scenarios that read
        the inbox (<code>process_inbox</code>, <code>validate_invoice</code>,{" "}
        <code>route_rfi</code>) will pick it up on their next run.
      </p>

      <div className="field">
        <div className="label">Template</div>
        <div className="input">
          <select
            value={tplId}
            onChange={(e) => onPickTemplate(e.target.value)}
          >
            {EMAIL_TEMPLATES.map((t) => (
              <option key={t.id} value={t.id}>
                {t.label}
              </option>
            ))}
          </select>
        </div>
        <div className="help">{tpl.description}</div>
        <div className="mock-mail-scenarios">
          <span className="lbl">Scenarios that use it:</span>
          {tpl.scenarios.map((s) => (
            <span key={s} className="chip-tag" title={`scenario: ${s}`}>
              {humanize(s)}
            </span>
          ))}
        </div>
      </div>

      <div className="row2">
        <div className="field">
          <div className="label">From</div>
          <div className="input mono">
            <input
              value={fields.from}
              onChange={(e) => update("from", e.target.value)}
            />
          </div>
        </div>
        <div className="field">
          <div className="label">From name</div>
          <div className="input">
            <input
              value={fields.from_name ?? ""}
              onChange={(e) => update("from_name", e.target.value)}
              placeholder="(optional display name)"
            />
          </div>
        </div>
      </div>

      <div className="field">
        <div className="label">Subject</div>
        <div className="input">
          <input
            value={fields.subject}
            onChange={(e) => update("subject", e.target.value)}
          />
        </div>
      </div>

      <div className="field">
        <div className="label">Body</div>
        <textarea
          className="input"
          rows={8}
          value={fields.body}
          onChange={(e) => update("body", e.target.value)}
        />
      </div>

      <div className="row2">
        <div className="field">
          <div className="label">Category</div>
          <div className="input">
            <select
              value={fields.category}
              onChange={(e) =>
                update("category", e.target.value as EmailTemplate["payload"]["category"])
              }
            >
              {["vendor", "applicant", "client", "internal"].map((c) => (
                <option key={c} value={c}>
                  {humanize(c)}
                </option>
              ))}
            </select>
          </div>
        </div>
        <div className="field">
          <div className="label">Trust status</div>
          <div className="input">
            <select
              value={fields.trust_status ?? ""}
              onChange={(e) =>
                update(
                  "trust_status",
                  (e.target.value || undefined) as
                    | "trusted"
                    | "impersonation"
                    | "new_contact"
                    | undefined
                )
              }
            >
              <option value="">—</option>
              <option value="trusted">Trusted</option>
              <option value="new_contact">New contact</option>
              <option value="impersonation">Impersonation</option>
            </select>
          </div>
        </div>
      </div>

      {fields.trust_status === "impersonation" && (
        <div className="field">
          <div className="label">Trust reason</div>
          <div className="input">
            <input
              value={fields.trust_reason ?? ""}
              onChange={(e) => update("trust_reason", e.target.value)}
              placeholder="why is it flagged"
            />
          </div>
        </div>
      )}
      {fields.category === "vendor" && (
        <div className="field">
          <div className="label">Vendor id (for validation)</div>
          <div className="input mono">
            <input
              value={fields.vendor_id ?? ""}
              onChange={(e) => update("vendor_id", e.target.value)}
              placeholder="vendor_xyz"
            />
          </div>
        </div>
      )}

      {flash && (
        <div className="test-flash">
          <span className="dot" />
          {flash}
        </div>
      )}

      <div className="lab-actions">
        <button
          className="btn btn--primary btn--sm"
          disabled={busy}
          onClick={onSend}
        >
          {busy ? "Sending…" : "Send to inbox"}
        </button>
        <button
          className="btn btn--ghost btn--sm"
          disabled={busy}
          onClick={() => setOverlay({})}
          title="Discard local edits, restore template defaults"
        >
          Reset to template
        </button>
      </div>

      <div className="lab-hint">
        Saved as <code>projects/inbox/mock_{tpl.id}_&lt;ts&gt;.json</code>.
        Remove by hand to clean up.
      </div>
    </section>
  );
}

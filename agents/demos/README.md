# Industry Demo Configurations

These are example agent team configurations for different industries.
To use one, copy the agents to your project's `.agent-house/agents/` directory.

## EPC (Engineering, Procurement, Construction)

| Agent | Mode | Role |
|-------|------|------|
| Project Director | oneshot | Scope evaluation, go/no-go decisions |
| Cost Controller | session | Cost estimation, SAP queries, budget tracking |
| Procurement Lead | session | Vendor management, bid evaluation |
| Safety Officer | oneshot | Compliance review, risk assessment |
| Proposal Writer | session | Technical proposal assembly |

**Use case:** Construction bid management, project estimation

## Hospital

| Agent | Mode | Role |
|-------|------|------|
| Duty Scheduler | session | Roster optimization, availability management |
| Healthcare Analyst | session | Patient volume analysis, predictions |
| Compliance Officer | oneshot | Labor law verification, regulatory checks |
| Report Generator | session | Monthly reports, dashboards |
| Staff Notifier | oneshot | Schedule distribution, communications |

**Use case:** Monthly roster generation, operational reporting

## HR

| Agent | Mode | Role |
|-------|------|------|
| Recruiter | session | Job descriptions, candidate sourcing |
| Resume Screener | oneshot | Application evaluation, shortlisting |
| Interview Coordinator | oneshot | Scheduling, logistics |
| Onboarding Specialist | session | Welcome packages, IT setup |

**Use case:** Hiring pipeline automation

## How to Use

```bash
# Copy EPC agents to your project
cp -r agents/demos/epc/* projects/my-epc-project/.agent-house/agents/

# Or configure team.json to use them
cat > projects/my-epc-project/.agent-house/team.json << EOF
{
  "team": [
    {"role": "project_director"},
    {"role": "cost_controller"},
    {"role": "procurement_lead"},
    {"role": "safety_officer"},
    {"role": "proposal_writer"}
  ],
  "pipeline": ["research", "planning", "review", "execution"]
}
EOF
```

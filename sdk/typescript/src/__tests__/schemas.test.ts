import { describe, it, expect } from "vitest";
import {
  // Schemas
  ModelSchema,
  AgentSchema,
  StepSchema,
  WorkflowSchema,
  TeamSchema,
  PlatformSchema,
  PrioritySchema,
  TargetSchema,
  DeploymentSchema,
  ClaudeCodeConfigSchema,
  KiroCLIConfigSchema,
  AWSAgentCoreConfigSchema,
  KubernetesConfigSchema,
  AgentKitLocalConfigSchema,
  // Utility functions
  parseAgent,
  safeParseAgent,
  parseTeam,
  safeParseTeam,
  parseDeployment,
  safeParseDeployment,
  // Mappings
  claudeCodeModels,
  kiroCliModels,
  bedrockModels,
  kiroCliTools,
} from "../index";

// =============================================================================
// Model Schema Tests
// =============================================================================

describe("ModelSchema", () => {
  it("should accept valid model names", () => {
    expect(ModelSchema.parse("haiku")).toBe("haiku");
    expect(ModelSchema.parse("sonnet")).toBe("sonnet");
    expect(ModelSchema.parse("opus")).toBe("opus");
  });

  it("should reject invalid model names", () => {
    expect(() => ModelSchema.parse("gpt-4")).toThrow();
    expect(() => ModelSchema.parse("")).toThrow();
  });

  it("should apply default", () => {
    // ModelSchema has a default of "sonnet"
    const schema = ModelSchema;
    expect(schema._def.defaultValue()).toBe("sonnet");
  });
});

// =============================================================================
// Agent Schema Tests
// =============================================================================

describe("AgentSchema", () => {
  it("should parse valid agent with all fields", () => {
    const agent = AgentSchema.parse({
      name: "test-agent",
      description: "A test agent",
      model: "sonnet",
      tools: ["Read", "Write"],
      skills: ["skill1"],
      dependencies: ["other-agent"],
      instructions: "You are a test agent.",
    });

    expect(agent.name).toBe("test-agent");
    expect(agent.description).toBe("A test agent");
    expect(agent.model).toBe("sonnet");
    expect(agent.tools).toEqual(["Read", "Write"]);
    expect(agent.skills).toEqual(["skill1"]);
    expect(agent.dependencies).toEqual(["other-agent"]);
    expect(agent.instructions).toBe("You are a test agent.");
  });

  it("should parse minimal agent with just name", () => {
    const agent = AgentSchema.parse({
      name: "minimal-agent",
    });

    expect(agent.name).toBe("minimal-agent");
    expect(agent.model).toBe("sonnet"); // default applied
  });

  it("should reject missing name", () => {
    expect(() => AgentSchema.parse({})).toThrow();
    expect(() => AgentSchema.parse({ description: "test" })).toThrow();
  });

  it("should support new fields: namespace, icon, requires, tasks", () => {
    const agent = AgentSchema.parse({
      name: "full-agent",
      namespace: "my-org",
      icon: "robot",
      requires: ["dep1", "dep2"],
      tasks: [
        { id: "task1", description: "First task", type: "command" },
        { id: "task2", type: "manual" },
      ],
    });

    expect(agent.namespace).toBe("my-org");
    expect(agent.icon).toBe("robot");
    expect(agent.requires).toEqual(["dep1", "dep2"]);
    expect(agent.tasks).toHaveLength(2);
  });
});

// =============================================================================
// Step Schema Tests
// =============================================================================

describe("StepSchema", () => {
  it("should parse valid step with all fields", () => {
    const step = StepSchema.parse({
      name: "research",
      agent: "researcher",
      depends_on: ["init"],
      inputs: [{ name: "topic", type: "string" }],
      outputs: [{ name: "results", type: "object" }],
    });

    expect(step.name).toBe("research");
    expect(step.agent).toBe("researcher");
    expect(step.depends_on).toEqual(["init"]);
    expect(step.inputs).toHaveLength(1);
    expect(step.outputs).toHaveLength(1);
  });

  it("should parse minimal step", () => {
    const step = StepSchema.parse({
      name: "step1",
      agent: "agent1",
    });

    expect(step.name).toBe("step1");
    expect(step.agent).toBe("agent1");
    expect(step.depends_on).toBeUndefined();
  });
});

// =============================================================================
// Workflow Schema Tests
// =============================================================================

describe("WorkflowSchema", () => {
  it("should parse valid workflow types", () => {
    for (const type of ["sequential", "parallel", "dag", "orchestrated"]) {
      const workflow = WorkflowSchema.parse({ type });
      expect(workflow.type).toBe(type);
    }
  });

  it("should apply default workflow type", () => {
    const workflow = WorkflowSchema.parse({});
    expect(workflow.type).toBe("orchestrated");
  });

  it("should parse workflow with steps", () => {
    const workflow = WorkflowSchema.parse({
      type: "dag",
      steps: [
        { name: "step1", agent: "agent1" },
        { name: "step2", agent: "agent2", depends_on: ["step1"] },
      ],
    });

    expect(workflow.steps).toHaveLength(2);
    expect(workflow.steps![1].depends_on).toEqual(["step1"]);
  });
});

// =============================================================================
// Team Schema Tests
// =============================================================================

describe("TeamSchema", () => {
  it("should parse valid team with all fields", () => {
    const team = TeamSchema.parse({
      name: "test-team",
      version: "1.0.0",
      description: "A test team",
      agents: ["agent1", "agent2"],
      orchestrator: "agent1",
      workflow: { type: "orchestrated" },
      context: "Shared context",
    });

    expect(team.name).toBe("test-team");
    expect(team.version).toBe("1.0.0");
    expect(team.agents).toEqual(["agent1", "agent2"]);
    expect(team.orchestrator).toBe("agent1");
  });

  it("should parse minimal team", () => {
    const team = TeamSchema.parse({
      name: "test-team",
      version: "1.0.0",
      agents: ["agent1"],
    });

    expect(team.name).toBe("test-team");
    expect(team.version).toBe("1.0.0");
    expect(team.agents).toEqual(["agent1"]);
  });

  it("should reject missing required fields", () => {
    expect(() =>
      TeamSchema.parse({
        name: "team",
        version: "1.0.0",
        // missing agents
      })
    ).toThrow();

    expect(() =>
      TeamSchema.parse({
        name: "team",
        agents: ["a1"],
        // missing version
      })
    ).toThrow();
  });
});

// =============================================================================
// Platform Schema Tests
// =============================================================================

describe("PlatformSchema", () => {
  it("should accept all valid platforms", () => {
    const platforms = [
      "claude-code",
      "gemini-cli",
      "kiro-cli",
      "adk-go",
      "crewai",
      "autogen",
      "aws-agentcore",
      "aws-eks",
      "azure-aks",
      "gcp-gke",
      "kubernetes",
      "docker-compose",
      "agentkit-local",
    ];
    for (const platform of platforms) {
      expect(PlatformSchema.parse(platform)).toBe(platform);
    }
  });

  it("should reject invalid platforms", () => {
    expect(() => PlatformSchema.parse("invalid")).toThrow();
  });
});

// =============================================================================
// Priority Schema Tests
// =============================================================================

describe("PrioritySchema", () => {
  it("should accept valid priorities", () => {
    expect(PrioritySchema.parse("p1")).toBe("p1");
    expect(PrioritySchema.parse("p2")).toBe("p2");
    expect(PrioritySchema.parse("p3")).toBe("p3");
  });

  it("should reject invalid priorities", () => {
    expect(() => PrioritySchema.parse("p0")).toThrow();
    expect(() => PrioritySchema.parse("high")).toThrow();
  });

  it("should have default of p2", () => {
    expect(PrioritySchema._def.defaultValue()).toBe("p2");
  });
});

// =============================================================================
// Platform Config Schema Tests
// =============================================================================

describe("ClaudeCodeConfigSchema", () => {
  it("should parse required fields", () => {
    const config = ClaudeCodeConfigSchema.parse({
      agentDir: ".claude/agents",
      format: "markdown",
    });
    expect(config.agentDir).toBe(".claude/agents");
    expect(config.format).toBe("markdown");
  });

  it("should reject missing required fields", () => {
    expect(() => ClaudeCodeConfigSchema.parse({})).toThrow();
    expect(() => ClaudeCodeConfigSchema.parse({ agentDir: "test" })).toThrow();
  });
});

describe("KiroCLIConfigSchema", () => {
  it("should parse required fields", () => {
    const config = KiroCLIConfigSchema.parse({
      pluginDir: "plugins/kiro/agents",
      format: "json",
    });
    expect(config.pluginDir).toBe("plugins/kiro/agents");
    expect(config.format).toBe("json");
  });
});

describe("AWSAgentCoreConfigSchema", () => {
  it("should parse required fields", () => {
    const config = AWSAgentCoreConfigSchema.parse({
      region: "us-east-1",
      foundationModel: "anthropic.claude-3-sonnet-20240229-v1:0",
      iac: "cdk",
      lambdaRuntime: "python3.11",
    });
    expect(config.region).toBe("us-east-1");
    expect(config.iac).toBe("cdk");
  });

  it("should reject missing required fields", () => {
    expect(() => AWSAgentCoreConfigSchema.parse({})).toThrow();
    expect(() =>
      AWSAgentCoreConfigSchema.parse({ region: "us-east-1" })
    ).toThrow();
  });
});

describe("KubernetesConfigSchema", () => {
  it("should parse required fields", () => {
    const config = KubernetesConfigSchema.parse({
      namespace: "multi-agent",
      helmChart: true,
    });
    expect(config.namespace).toBe("multi-agent");
    expect(config.helmChart).toBe(true);
  });

  it("should parse with resource limits", () => {
    const config = KubernetesConfigSchema.parse({
      namespace: "test",
      helmChart: false,
      resourceLimits: { cpu: "1000m", memory: "1Gi" },
    });
    expect(config.resourceLimits?.cpu).toBe("1000m");
    expect(config.resourceLimits?.memory).toBe("1Gi");
  });
});

describe("AgentKitLocalConfigSchema", () => {
  it("should parse required transport", () => {
    const config = AgentKitLocalConfigSchema.parse({
      transport: "stdio",
    });
    expect(config.transport).toBe("stdio");
    expect(config.port).toBeUndefined();
  });

  it("should parse http transport with port", () => {
    const config = AgentKitLocalConfigSchema.parse({
      transport: "http",
      port: 8080,
    });
    expect(config.transport).toBe("http");
    expect(config.port).toBe(8080);
  });
});

// =============================================================================
// Target Schema Tests
// =============================================================================

describe("TargetSchema", () => {
  it("should parse valid target", () => {
    const target = TargetSchema.parse({
      name: "local-claude",
      platform: "claude-code",
      priority: "p1",
      output: ".claude/agents",
    });

    expect(target.name).toBe("local-claude");
    expect(target.platform).toBe("claude-code");
    expect(target.priority).toBe("p1");
    expect(target.output).toBe(".claude/agents");
  });

  it("should apply default priority", () => {
    const target = TargetSchema.parse({
      name: "test",
      platform: "claude-code",
    });
    expect(target.priority).toBe("p2");
  });

  it("should support platform-specific configs", () => {
    const target = TargetSchema.parse({
      name: "claude-local",
      platform: "claude-code",
      claudeCode: {
        agentDir: ".claude/agents",
        format: "markdown",
      },
    });
    expect(target.claudeCode?.agentDir).toBe(".claude/agents");
  });
});

// =============================================================================
// Deployment Schema Tests
// =============================================================================

describe("DeploymentSchema", () => {
  it("should parse valid deployment", () => {
    const deployment = DeploymentSchema.parse({
      team: "test-team",
      targets: [
        {
          name: "local",
          platform: "claude-code",
        },
      ],
    });

    expect(deployment.team).toBe("test-team");
    expect(deployment.targets).toHaveLength(1);
  });

  it("should reject missing targets", () => {
    expect(() =>
      DeploymentSchema.parse({
        team: "test",
      })
    ).toThrow();
  });

  it("should accept $schema field", () => {
    const deployment = DeploymentSchema.parse({
      $schema: "../schema/deployment.schema.json",
      team: "test",
      targets: [{ name: "t", platform: "claude-code" }],
    });
    expect(deployment.$schema).toBe("../schema/deployment.schema.json");
  });
});

// =============================================================================
// Utility Function Tests
// =============================================================================

describe("parseAgent", () => {
  it("should parse valid agent", () => {
    const agent = parseAgent({
      name: "test",
    });
    expect(agent.name).toBe("test");
  });

  it("should throw on invalid agent", () => {
    expect(() => parseAgent({})).toThrow();
  });
});

describe("safeParseAgent", () => {
  it("should return success for valid agent", () => {
    const result = safeParseAgent({
      name: "test",
    });
    expect(result.success).toBe(true);
    if (result.success) {
      expect(result.data.name).toBe("test");
    }
  });

  it("should return error for invalid agent", () => {
    const result = safeParseAgent({});
    expect(result.success).toBe(false);
  });
});

describe("parseTeam", () => {
  it("should parse valid team", () => {
    const team = parseTeam({
      name: "test",
      version: "1.0.0",
      agents: ["agent1"],
    });
    expect(team.name).toBe("test");
  });

  it("should throw on invalid team", () => {
    expect(() => parseTeam({ name: "test" })).toThrow();
  });
});

describe("safeParseTeam", () => {
  it("should return success for valid team", () => {
    const result = safeParseTeam({
      name: "test",
      version: "1.0.0",
      agents: ["agent1"],
    });
    expect(result.success).toBe(true);
  });

  it("should return error for invalid team", () => {
    const result = safeParseTeam({ name: "test" });
    expect(result.success).toBe(false);
  });
});

describe("parseDeployment", () => {
  it("should parse valid deployment", () => {
    const deployment = parseDeployment({
      team: "test",
      targets: [{ name: "t", platform: "claude-code" }],
    });
    expect(deployment.team).toBe("test");
  });

  it("should throw on invalid deployment", () => {
    expect(() => parseDeployment({ team: "test" })).toThrow();
  });
});

describe("safeParseDeployment", () => {
  it("should return success for valid deployment", () => {
    const result = safeParseDeployment({
      team: "test",
      targets: [{ name: "t", platform: "claude-code" }],
    });
    expect(result.success).toBe(true);
  });

  it("should return error for invalid deployment", () => {
    const result = safeParseDeployment({ team: "test" });
    expect(result.success).toBe(false);
  });
});

// =============================================================================
// Mapping Tests
// =============================================================================

describe("claudeCodeModels", () => {
  it("should have all model mappings", () => {
    expect(claudeCodeModels.haiku).toBe("haiku");
    expect(claudeCodeModels.sonnet).toBe("sonnet");
    expect(claudeCodeModels.opus).toBe("opus");
  });
});

describe("kiroCliModels", () => {
  it("should have all model mappings", () => {
    expect(kiroCliModels.haiku).toBe("claude-haiku-35");
    expect(kiroCliModels.sonnet).toBe("claude-sonnet-4");
    expect(kiroCliModels.opus).toBe("claude-opus-4");
  });
});

describe("bedrockModels", () => {
  it("should have all model mappings", () => {
    expect(bedrockModels.haiku).toContain("haiku");
    expect(bedrockModels.sonnet).toContain("sonnet");
    expect(bedrockModels.opus).toContain("opus");
  });
});

describe("kiroCliTools", () => {
  it("should have all tool mappings", () => {
    expect(kiroCliTools.WebSearch).toBe("web_search");
    expect(kiroCliTools.WebFetch).toBe("web_fetch");
    expect(kiroCliTools.Read).toBe("read");
    expect(kiroCliTools.Write).toBe("write");
    expect(kiroCliTools.Glob).toBe("glob");
    expect(kiroCliTools.Grep).toBe("grep");
    expect(kiroCliTools.Bash).toBe("bash");
    expect(kiroCliTools.Edit).toBe("edit");
    expect(kiroCliTools.Task).toBe("task");
  });
});

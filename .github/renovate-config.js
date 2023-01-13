module.exports = {
  platform: 'github',
  extends: ['config:base', ':disableDependencyDashboard', ":semanticCommitTypeAll(deps)"],
  packageRules: [
    {
      matchUpdateTypes: ["minor", "patch"],
      automerge: true
    },
    // Istio has to be considered as a whole
    {
      matchPaths: ["addons/istio/**"],
      groupName: "Istio Helm Chart"
    }
  ],
  username: "nearform-renovate-app",
  gitAuthor: '"Nearform" <nearform@nearform.com>',
  automergeType: "branch",
  platformAutomerge: true,
  repositories: ['nearform/k8s-kurated-addons'],
  enabledManagers: ['asdf', 'helmv3', 'github-actions'],
  prConcurrentLimit: 0,
  branchConcurrentLimit: 0,
  branchNameStrict: true,
  baseBranches: ['main'],
  dependencyDashboard: false,
  rebaseWhen : "auto"
};

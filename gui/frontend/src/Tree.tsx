// Tree.tsx — the vault file tree in the left pane. Folds folders; a click on a
// note opens it. The flat TreeEntry list from the Go side (already sorted by
// path) is nested here for display.

import { useMemo, useState } from "react";
import type { TreeEntry } from "./api";

interface TreeNode {
  path: string;
  name: string;
  dir: boolean;
  children: TreeNode[];
}

function build(entries: TreeEntry[]): TreeNode[] {
  const roots: TreeNode[] = [];
  const byPath = new Map<string, TreeNode>();
  for (const e of entries) {
    const node: TreeNode = { ...e, children: [] };
    byPath.set(e.path, node);
    const slash = e.path.lastIndexOf("/");
    const parent = slash < 0 ? null : byPath.get(e.path.slice(0, slash));
    if (parent) parent.children.push(node);
    else roots.push(node);
  }
  return roots;
}

interface TreeProps {
  entries: TreeEntry[];
  current: string;
  onOpen(path: string): void;
}

export default function Tree({ entries, current, onOpen }: TreeProps) {
  const roots = useMemo(() => build(entries), [entries]);
  const [collapsed, setCollapsed] = useState<Set<string>>(new Set());

  function toggle(path: string) {
    setCollapsed((c) => {
      const n = new Set(c);
      n.has(path) ? n.delete(path) : n.add(path);
      return n;
    });
  }

  function render(nodes: TreeNode[], depth: number): React.ReactNode {
    return nodes.map((n) => {
      const pad = { paddingLeft: 6 + depth * 14 };
      if (n.dir) {
        const isCollapsed = collapsed.has(n.path);
        return (
          <div key={n.path}>
            <div className="row dir" style={pad} onClick={() => toggle(n.path)}>
              <span className="tw">{isCollapsed ? "▸" : "▾"}</span>
              {n.name}
            </div>
            {!isCollapsed && render(n.children, depth + 1)}
          </div>
        );
      }
      return (
        <div
          key={n.path}
          className={"row note" + (n.path === current ? " active" : "")}
          style={pad}
          onClick={() => onOpen(n.path)}
        >
          {n.name}
        </div>
      );
    });
  }

  if (entries.length === 0) {
    return <div className="tree-empty">No notes yet.</div>;
  }
  return <div className="tree">{render(roots, 0)}</div>;
}

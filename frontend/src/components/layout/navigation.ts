import {
  LayoutDashboard,
  MessageSquare,
  Bot,
  FileText,
  Megaphone,
  Settings,
  Users,
  Contact,
  Workflow,
  Sparkles,
  Key,
  UserX,
  MessageSquareText,
  Webhook,
  BarChart3,
  ShieldCheck,
  Zap,
  Shield,
  LineChart,
  Tags,
  PhoneCall,
  PhoneForwarded,
  ScrollText,
  Palette
} from 'lucide-vue-next'
import type { Component } from 'vue'

export interface NavItem {
  name: string
  path: string
  icon: Component
  permission?: string
  childPermissions?: string[]
  children?: NavItem[]
}

export interface NavSection {
  label: string
  items: NavItem[]
  /** Permissions needed to show section — at least one must pass */
  permissions: string[]
  /** Pin to bottom of sidebar */
  pinBottom?: boolean
}

export const navigationSections: NavSection[] = [
  {
    label: 'nav.sectionMain',
    permissions: ['analytics', 'chat'],
    items: [
      {
        name: 'nav.dashboard',
        path: '/app/dashboard',
        icon: LayoutDashboard,
        permission: 'analytics'
      },
      {
        name: 'nav.chat',
        path: '/app/chat',
        icon: MessageSquare,
        permission: 'chat'
      },
    ]
  },
  {
    label: 'nav.sectionMessaging',
    permissions: ['settings.chatbot', 'chatbot.keywords', 'flows.chatbot', 'chatbot.ai', 'transfers', 'campaigns', 'templates', 'flows.whatsapp'],
    items: [
      {
        name: 'nav.chatbot',
        path: '/app/chatbot',
        icon: Bot,
        permission: 'settings.chatbot',
        childPermissions: ['settings.chatbot', 'chatbot.keywords', 'flows.chatbot', 'chatbot.ai', 'transfers'],
        children: [
          { name: 'nav.overview', path: '/app/chatbot', icon: Bot, permission: 'settings.chatbot' },
          { name: 'nav.keywords', path: '/app/chatbot/keywords', icon: Key, permission: 'chatbot.keywords' },
          { name: 'nav.flows', path: '/app/chatbot/flows', icon: Workflow, permission: 'flows.chatbot' },
          { name: 'nav.aiContexts', path: '/app/chatbot/ai', icon: Sparkles, permission: 'chatbot.ai' },
          { name: 'nav.transfers', path: '/app/chatbot/transfers', icon: UserX, permission: 'transfers' }
        ]
      },
      {
        name: 'nav.campaigns',
        path: '/app/campaigns',
        icon: Megaphone,
        permission: 'campaigns'
      },
      {
        name: 'nav.templates',
        path: '/app/templates',
        icon: FileText,
        permission: 'templates'
      },
      {
        name: 'nav.flows',
        path: '/app/flows',
        icon: Workflow,
        permission: 'flows.whatsapp'
      },
    ]
  },
  {
    label: 'nav.sectionCalling',
    permissions: ['call_logs', 'ivr_flows', 'call_transfers'],
    items: [
      { name: 'nav.callLogs', path: '/app/calling/logs', icon: PhoneCall, permission: 'call_logs' },
      { name: 'nav.ivrFlows', path: '/app/calling/ivr-flows', icon: Workflow, permission: 'ivr_flows' },
      { name: 'nav.callTransfers', path: '/app/calling/transfers', icon: PhoneForwarded, permission: 'call_transfers' },
    ]
  },
  {
    label: 'nav.sectionAnalytics',
    permissions: ['analytics.agents', 'analytics'],
    items: [
      {
        name: 'nav.agentAnalytics',
        path: '/app/analytics/agents',
        icon: BarChart3,
        permission: 'analytics.agents'
      },
      {
        name: 'nav.metaInsights',
        path: '/app/analytics/meta-insights',
        icon: LineChart,
        permission: 'analytics'
      },
    ]
  },
  {
    label: '',
    permissions: ['settings.general', 'settings.chatbot', 'accounts', 'contacts', 'canned_responses', 'tags', 'teams', 'users', 'roles', 'api_keys', 'webhooks', 'custom_actions', 'settings.sso', 'audit_logs'],
    pinBottom: true,
    items: [
      {
        name: 'nav.settings',
        path: '/app/settings',
        icon: Settings,
        permission: 'settings.general',
        childPermissions: ['settings.general', 'settings.chatbot', 'accounts', 'contacts', 'canned_responses', 'tags', 'teams', 'users', 'roles', 'api_keys', 'webhooks', 'custom_actions', 'settings.sso', 'audit_logs'],
        children: [
          { name: 'nav.general', path: '/app/settings', icon: Settings, permission: 'settings.general' },
          { name: 'nav.chatbot', path: '/app/settings/chatbot', icon: Bot, permission: 'settings.chatbot' },
          { name: 'nav.accounts', path: '/app/settings/accounts', icon: Users, permission: 'accounts' },
          { name: 'nav.contacts', path: '/app/settings/contacts', icon: Contact, permission: 'contacts' },
          { name: 'nav.cannedResponses', path: '/app/settings/canned-responses', icon: MessageSquareText, permission: 'canned_responses' },
          { name: 'nav.tags', path: '/app/settings/tags', icon: Tags, permission: 'tags' },
          { name: 'nav.teams', path: '/app/settings/teams', icon: Users, permission: 'teams' },
          { name: 'nav.users', path: '/app/settings/users', icon: Users, permission: 'users' },
          { name: 'nav.roles', path: '/app/settings/roles', icon: Shield, permission: 'roles' },
          { name: 'nav.apiKeys', path: '/app/settings/api-keys', icon: Key, permission: 'api_keys' },
          { name: 'nav.webhooks', path: '/app/settings/webhooks', icon: Webhook, permission: 'webhooks' },
          { name: 'nav.customActions', path: '/app/settings/custom-actions', icon: Zap, permission: 'custom_actions' },
          { name: 'nav.sso', path: '/app/settings/sso', icon: ShieldCheck, permission: 'settings.sso' },
          { name: 'nav.auditLogs', path: '/app/settings/audit-logs', icon: ScrollText, permission: 'audit_logs' },
          { name: 'nav.branding', path: '/app/settings/branding', icon: Palette, permission: 'settings.general' }
        ]
      }
    ]
  }
]

// Flat list for backward compatibility (used by AppLayout computed)
export const navigationItems: NavItem[] = navigationSections.flatMap(s => s.items)

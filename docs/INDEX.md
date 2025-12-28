# Usual Store Documentation Index

Welcome to the Usual Store documentation! This index will help you find the information you need quickly.

## 📚 Documentation Structure

```
docs/
├── INDEX.md                              ← You are here
├── TERRAFORM-INFRASTRUCTURE.md           ← Terraform overview
│
├── setup/                                ← Getting started guides
├── guides/                               ← Feature & implementation guides
├── summaries/                            ← Project summaries
└── ipv6-docker-setup/                    ← IPv6 configuration archive
```

## 🚀 Getting Started

### Initial Setup
| Document | Description |
|----------|-------------|
| [AUTHENTICATION-SETUP.md](setup/AUTHENTICATION-SETUP.md) | Authentication system setup |
| [HOW-TO-ACCESS.md](setup/HOW-TO-ACCESS.md) | Access instructions for services |
| [TYPESCRIPT-FRONTEND-SETUP.md](setup/TYPESCRIPT-FRONTEND-SETUP.md) | TypeScript frontend setup |
| [MATERIAL-UI-SETUP.md](setup/MATERIAL-UI-SETUP.md) | Material-UI configuration |
| [REACT-FRONTEND-SETUP.md](setup/REACT-FRONTEND-SETUP.md) | React frontend setup |
| [STRIPE-SETUP.md](setup/STRIPE-SETUP.md) | Stripe payment integration |
| [LOGIN-CREDENTIALS.md](setup/LOGIN-CREDENTIALS.md) | Default login credentials |

## 📖 Feature Guides

### Infrastructure & DevOps
| Document | Description |
|----------|-------------|
| [TERRAFORM-INFRASTRUCTURE.md](TERRAFORM-INFRASTRUCTURE.md) | Terraform infrastructure overview |
| [KUBERNETES-README.md](guides/KUBERNETES-README.md) | Kubernetes deployment guide |
| [KUBERNETES-DASHBOARD-GUIDE.md](guides/KUBERNETES-DASHBOARD-GUIDE.md) | Kubernetes dashboard & monitoring |
| [AWS-LAMBDA-GUIDE.md](../terraform-tenants/AWS-LAMBDA-GUIDE.md) | **AWS Lambda + RDS deployment** ⭐ |
| [MONGODB-GUIDE.md](../terraform-tenants/MONGODB-GUIDE.md) | **MongoDB option** (local & AWS) 🍃 |
| [IPv6-SUCCESS.md](guides/IPv6-SUCCESS.md) | IPv6 implementation |
| [GITHUB-ACTIONS-CACHE-FIX.md](guides/GITHUB-ACTIONS-CACHE-FIX.md) | CI/CD cache optimization |

### Messaging & Communication
| Document | Description |
|----------|-------------|
| [MESSAGING-SYSTEM-ARCHITECTURE.md](guides/MESSAGING-SYSTEM-ARCHITECTURE.md) | Complete messaging architecture |
| [KAFKA-MESSAGING-IMPLEMENTATION.md](guides/KAFKA-MESSAGING-IMPLEMENTATION.md) | Kafka setup guide |
| [KAFKA-TESTING-GUIDE.md](guides/KAFKA-TESTING-GUIDE.md) | Kafka testing instructions |
| [KAFKA-MESSAGING-ARCHITECTURE.md](guides/KAFKA-MESSAGING-ARCHITECTURE.md) | Kafka architecture details |
| [EMAIL-QUICK-REFERENCE.md](guides/EMAIL-QUICK-REFERENCE.md) | Email system reference |
| [USER-CREATION-EMAIL-NOTIFICATIONS.md](guides/USER-CREATION-EMAIL-NOTIFICATIONS.md) | Email notifications for new users |

### Support System
| Document | Description |
|----------|-------------|
| [LIVE-SUPPORT-CHAT.md](guides/LIVE-SUPPORT-CHAT.md) | Live support implementation |
| [LIVE-SUPPORT-QUICKSTART.md](guides/LIVE-SUPPORT-QUICKSTART.md) | Quick start for support system |

### User Management
| Document | Description |
|----------|-------------|
| [USER-MANAGEMENT-FEATURES.md](guides/USER-MANAGEMENT-FEATURES.md) | User management system |
| [SERVER-SIDE-SORTING-SEARCH.md](guides/SERVER-SIDE-SORTING-SEARCH.md) | Server-side data operations |

### Frontend Applications
| Document | Description |
|----------|-------------|
| [REACT-REDUX-FRONTEND.md](guides/REACT-REDUX-FRONTEND.md) | Redux frontend implementation |
| [AI-ASSISTANT-README.md](guides/AI-ASSISTANT-README.md) | AI assistant integration |

### Multi-Tenancy (Terraform-Based)
| Document | Description |
|----------|-------------|
| [SIMPLIFIED-ARCHITECTURE.md](../terraform-tenants/SIMPLIFIED-ARCHITECTURE.md) | NO master DB - complete separation |
| [terraform-tenants/README.md](../terraform-tenants/README.md) | Interactive script to add customers |
| [MULTI-TENANT-ARCHITECTURE.md](../MULTI-TENANT-ARCHITECTURE.md) | Architecture overview & responsibilities |

### Development & Testing
| Document | Description |
|----------|-------------|
| [DEVELOPMENT-WORKFLOW.md](guides/DEVELOPMENT-WORKFLOW.md) | How to update code changes |
| [TABLE-DRIVEN-TESTS.md](guides/TABLE-DRIVEN-TESTS.md) | Testing methodology |
| [TEST-SUMMARY.md](guides/TEST-SUMMARY.md) | Testing overview |
| [AI-TEST-COVERAGE-SUMMARY.md](guides/AI-TEST-COVERAGE-SUMMARY.md) | Test coverage analysis |
| [LINTER-FIXES.md](guides/LINTER-FIXES.md) | Code quality improvements |
| [GO-PACKAGE-CONFLICT-FIX.md](guides/GO-PACKAGE-CONFLICT-FIX.md) | Package conflict resolution |

### Observability & Monitoring
| Document | Description |
|----------|-------------|
| [OPENTELEMETRY-TRACING-GUIDE.md](guides/OPENTELEMETRY-TRACING-GUIDE.md) | Distributed tracing setup |

### Libraries & Dependencies
| Document | Description |
|----------|-------------|
| [GO-LIBRARIES-EVALUATION-2025.md](guides/GO-LIBRARIES-EVALUATION-2025.md) | Go libraries evaluation |
| [IMPLEMENTATION-COMPLETE-2025-LIBRARIES.md](guides/IMPLEMENTATION-COMPLETE-2025-LIBRARIES.md) | Libraries implementation status |

### Tools & Commands
| Document | Description |
|----------|-------------|
| [MAKEFILE-COMMANDS.md](guides/MAKEFILE-COMMANDS.md) | Available Make commands |
| [MAKEFILE-QUICK-REFERENCE.md](guides/MAKEFILE-QUICK-REFERENCE.md) | Quick Makefile reference |
| [DOCUMENTATION.md](guides/DOCUMENTATION.md) | General documentation guide |

## 🏗️ Terraform Documentation

Terraform-specific documentation is located in `../terraform/`:

| Document | Description |
|----------|-------------|
| [terraform/README.md](../terraform/README.md) | Main Terraform guide |
| [terraform/POLICY-EXAMPLES.md](../terraform/POLICY-EXAMPLES.md) | OPA policy examples |
| [terraform/QUICK-START.md](../terraform/QUICK-START.md) | Quick start guide |
| [terraform/CHEAT-SHEET.md](../terraform/CHEAT-SHEET.md) | Command reference |
| [TERRAFORM-SETUP-COMPLETE.md](guides/TERRAFORM-SETUP-COMPLETE.md) | Setup completion guide |

## 📊 Project Summaries

| Document | Description |
|----------|-------------|
| [SESSION-COMPLETE-SUMMARY.md](summaries/SESSION-COMPLETE-SUMMARY.md) | Complete project summary |

## 🌐 IPv6 Documentation (Archive)

Historical IPv6 setup documentation is available in [ipv6-docker-setup/](ipv6-docker-setup/):

- [IPv6-SETUP.md](ipv6-docker-setup/IPv6-SETUP.md) - Initial setup guide
- [IPv6-SOLUTIONS.md](ipv6-docker-setup/IPv6-SOLUTIONS.md) - Solutions to common issues
- [QUICKSTART-IPv6.md](ipv6-docker-setup/QUICKSTART-IPv6.md) - Quick start guide
- And more...

## 🎯 Quick Links by Role

### For Developers
- [AUTHENTICATION-SETUP.md](setup/AUTHENTICATION-SETUP.md) - Auth setup
- [TABLE-DRIVEN-TESTS.md](guides/TABLE-DRIVEN-TESTS.md) - Testing
- [MAKEFILE-COMMANDS.md](guides/MAKEFILE-COMMANDS.md) - Commands
- [GO-LIBRARIES-EVALUATION-2025.md](guides/GO-LIBRARIES-EVALUATION-2025.md) - Libraries

### For DevOps
- [TERRAFORM-INFRASTRUCTURE.md](TERRAFORM-INFRASTRUCTURE.md) - Infrastructure
- [KUBERNETES-README.md](guides/KUBERNETES-README.md) - K8s deployment
- [KAFKA-MESSAGING-IMPLEMENTATION.md](guides/KAFKA-MESSAGING-IMPLEMENTATION.md) - Kafka
- [OPENTELEMETRY-TRACING-GUIDE.md](guides/OPENTELEMETRY-TRACING-GUIDE.md) - Monitoring

### For Frontend Developers
- [REACT-FRONTEND-SETUP.md](setup/REACT-FRONTEND-SETUP.md) - React setup
- [TYPESCRIPT-FRONTEND-SETUP.md](setup/TYPESCRIPT-FRONTEND-SETUP.md) - TypeScript setup
- [MATERIAL-UI-SETUP.md](setup/MATERIAL-UI-SETUP.md) - UI components
- [REACT-REDUX-FRONTEND.md](guides/REACT-REDUX-FRONTEND.md) - Redux

### For System Administrators
- [HOW-TO-ACCESS.md](setup/HOW-TO-ACCESS.md) - Access guide
- [LOGIN-CREDENTIALS.md](setup/LOGIN-CREDENTIALS.md) - Credentials
- [USER-MANAGEMENT-FEATURES.md](guides/USER-MANAGEMENT-FEATURES.md) - User management

## 📝 Contributing

When adding new documentation:
1. Place setup guides in `setup/`
2. Place feature guides in `guides/`
3. Update this INDEX.md
4. Follow the naming convention: `TOPIC-NAME.md`

## 🔗 External Resources

- Main README: [../README.md](../README.md)
- Terraform README: [../terraform/README.md](../terraform/README.md)
- React Frontend: [../react-frontend/README.md](../react-frontend/README.md)
- TypeScript Frontend: [../typescript-frontend/README.md](../typescript-frontend/README.md)
- Redux Frontend: [../react-redux-frontend/README.md](../react-redux-frontend/README.md)
- Kubernetes: [../k8s/README.md](../k8s/README.md)

---

**Last Updated**: December 26, 2025  
**Total Documents**: 50+  
**Maintained By**: Usual Store Team


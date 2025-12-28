# Example: Multiple Tenants with Different Access Levels

postgres_password = "yourpassword"

tenants = {
  # Tenant 1: E-commerce store with multiple developers
  online_store = {
    database_name = "online_store_2025"
    tenant_name   = "Online Store LLC"
    plan          = "enterprise"
    
    admins = [
      {
        username = "store_owner"
        password = "OwnerSecure123!"
        email    = "owner@onlinestore.com"
      },
      {
        username = "store_manager"
        password = "ManagerSecure123!"
        email    = "manager@onlinestore.com"
      }
    ]
    
    developers = [
      {
        username = "backend_dev"
        password = "BackendDev123!"
        email    = "backend@onlinestore.com"
      },
      {
        username = "frontend_dev"
        password = "FrontendDev123!"
        email    = "frontend@onlinestore.com"
      },
      {
        username = "qa_engineer"
        password = "QAEngineer123!"
        email    = "qa@onlinestore.com"
      }
    ]
    
    customers = [
      {
        username = "reporting_user"
        password = "ReportView123!"
        email    = "reports@onlinestore.com"
      },
      {
        username = "analytics_user"
        password = "AnalyticsView123!"
        email    = "analytics@onlinestore.com"
      }
    ]
  }

  # Tenant 2: Small startup with minimal users
  startup_app = {
    database_name = "startup_db"
    tenant_name   = "Startup App"
    plan          = "free"
    
    admins = [
      {
        username = "startup_admin"
        password = "StartupAdmin123!"
        email    = "admin@startup.com"
      }
    ]
    
    developers = [
      {
        username = "fullstack_dev"
        password = "FullstackDev123!"
        email    = "dev@startup.com"
      }
    ]
    
    customers = []  # No customer users needed
  }

  # Tenant 3: Manufacturing company
  manufacturing = {
    database_name = "manufacturing_erp"
    tenant_name   = "Manufacturing Corp"
    plan          = "professional"
    
    admins = [
      {
        username = "it_admin"
        password = "ITAdmin123!"
        email    = "it@manufacturing.com"
      }
    ]
    
    developers = [
      {
        username = "erp_developer"
        password = "ERPDev123!"
        email    = "erp-dev@manufacturing.com"
      }
    ]
    
    customers = [
      {
        username = "warehouse_viewer"
        password = "WarehouseView123!"
        email    = "warehouse@manufacturing.com"
      },
      {
        username = "sales_viewer"
        password = "SalesView123!"
        email    = "sales@manufacturing.com"
      }
    ]
  }

  # Tenant 4: Agency managing multiple clients
  agency_platform = {
    database_name = "agency_clients"
    tenant_name   = "Digital Agency"
    plan          = "professional"
    
    admins = [
      {
        username = "agency_owner"
        password = "AgencyOwner123!"
        email    = "owner@agency.com"
      }
    ]
    
    developers = [
      {
        username = "agency_dev1"
        password = "AgencyDev1!"
        email    = "dev1@agency.com"
      },
      {
        username = "agency_dev2"
        password = "AgencyDev2!"
        email    = "dev2@agency.com"
      }
    ]
    
    customers = [
      {
        username = "client_a_viewer"
        password = "ClientAView123!"
        email    = "clienta@example.com"
      },
      {
        username = "client_b_viewer"
        password = "ClientBView123!"
        email    = "clientb@example.com"
      }
    ]
  }
}


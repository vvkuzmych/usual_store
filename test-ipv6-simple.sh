#!/bin/bash

# Simple IPv6 Test Script for Usual Store
# Run this anytime to verify IPv6 is working

echo "🧪 Testing IPv6 Configuration..."
echo ""

# Test 1: Network IPv6 Status
echo "1️⃣  Network IPv6 Enabled:"
docker network inspect usual_store_usualstore_network 2>/dev/null | grep "EnableIPv6" && echo "   ✅ IPv6 Enabled" || echo "   ❌ IPv6 Not Enabled"
echo ""

# Test 2: Container IPv6 Address
echo "2️⃣  Database IPv6 Address:"
IPV6_ADDR=$(docker inspect usual_store-database-1 2>/dev/null | grep '"GlobalIPv6Address"' | head -1 | awk -F'"' '{print $4}')
if [ -n "$IPV6_ADDR" ] && [ "$IPV6_ADDR" != "" ]; then
    echo "   ✅ Has IPv6: $IPV6_ADDR"
else
    echo "   ⚠️  No IPv6 address assigned"
fi
echo ""

# Test 3: Port Binding
echo "3️⃣  Port Bindings:"
docker compose ps 2>/dev/null | grep database | grep -o "\[::1\]:5433" && echo "   ✅ IPv6 port binding active" || echo "   ❌ No IPv6 port binding"
echo ""

# Test 4: Database Connection (IPv6)
echo "4️⃣  Database Connection via IPv6 ([::1]:5433):"
if psql "postgres://postgres:password@[::1]:5433/usualstore?sslmode=disable" -c "SELECT 1;" > /dev/null 2>&1; then
    echo "   ✅ IPv6 connection works!"
    psql "postgres://postgres:password@[::1]:5433/usualstore?sslmode=disable" -t -c "SELECT '   Server IP: ' || inet_server_addr();" 2>/dev/null
else
    echo "   ❌ IPv6 connection failed"
fi
echo ""

# Test 5: Database Connection (IPv4 - for comparison)
echo "5️⃣  Database Connection via IPv4 (127.0.0.1:5433):"
if psql "postgres://postgres:password@127.0.0.1:5433/usualstore?sslmode=disable" -c "SELECT 1;" > /dev/null 2>&1; then
    echo "   ✅ IPv4 connection works!"
    psql "postgres://postgres:password@127.0.0.1:5433/usualstore?sslmode=disable" -t -c "SELECT '   Server IP: ' || inet_server_addr();" 2>/dev/null
else
    echo "   ❌ IPv4 connection failed"
fi
echo ""

# Test 6: Web App
echo "6️⃣  Web Application:"
if curl -s -o /dev/null -w "" http://127.0.0.1:4000 2>/dev/null; then
    echo "   ✅ Web app accessible (IPv4)"
else
    echo "   ⚠️  Web app not accessible"
fi
echo ""

# Summary
echo "═══════════════════════════════════════"
echo "📊 SUMMARY:"
echo ""
if [ -n "$IPV6_ADDR" ] && [ "$IPV6_ADDR" != "" ]; then
    echo "✅ IPv6 is ENABLED and WORKING"
    echo ""
    echo "   Connect via IPv6:"
    echo "   psql postgres://postgres:password@[::1]:5433/usualstore"
    echo ""
    echo "   Connect via IPv4:"
    echo "   psql postgres://postgres:password@127.0.0.1:5433/usualstore"
else
    echo "⚠️  IPv6 configuration needs attention"
    echo ""
    echo "   Check: docker-compose.yml has enable_ipv6: true"
    echo "   Check: Docker Desktop Settings → Network → Dual IPv4/IPv6"
fi
echo "═══════════════════════════════════════"


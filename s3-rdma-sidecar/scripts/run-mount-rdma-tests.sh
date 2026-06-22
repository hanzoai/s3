#!/bin/bash

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
COMPOSE_FILE="docker-compose.mount-rdma.yml"
PROJECT_NAME="hanzo-rdma-mount"

# Function to show usage
show_usage() {
    echo -e "${BLUE}🚀 Hanzo RDMA Mount Test Runner${NC}"
    echo "===================================="
    echo ""
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  start           Start the RDMA mount environment"
    echo "  stop            Stop and cleanup the environment"
    echo "  restart         Restart the environment"
    echo "  status          Show status of all services"
    echo "  logs [service]  Show logs for all services or specific service"
    echo "  test            Run integration tests"
    echo "  perf            Run performance tests"
    echo "  shell           Open shell in mount container"
    echo "  cleanup         Full cleanup including volumes"
    echo ""
    echo "Services:"
    echo "  hanzo-master    Hanzo master server"
    echo "  hanzo-volume    Hanzo volume server"
    echo "  hanzo-filer     Hanzo filer server"
    echo "  rdma-engine         RDMA engine (Rust)"
    echo "  rdma-sidecar        RDMA sidecar (Go)"
    echo "  hanzo-mount     Hanzo mount with RDMA"
    echo ""
    echo "Examples:"
    echo "  $0 start                    # Start all services"
    echo "  $0 logs hanzo-mount     # Show mount logs"
    echo "  $0 test                     # Run integration tests"
    echo "  $0 perf                     # Run performance tests"
    echo "  $0 shell                    # Open shell in mount container"
}

# Function to check if Docker Compose is available
check_docker_compose() {
    if ! command -v docker-compose >/dev/null 2>&1 && ! docker compose version >/dev/null 2>&1; then
        echo -e "${RED}❌ Docker Compose is not available${NC}"
        echo "Please install Docker Compose to continue"
        exit 1
    fi
    
    # Use docker compose if available, otherwise docker-compose
    if docker compose version >/dev/null 2>&1; then
        DOCKER_COMPOSE="docker compose"
    else
        DOCKER_COMPOSE="docker-compose"
    fi
}

# Function to build required images
build_images() {
    echo -e "${BLUE}🔨 Building required Docker images...${NC}"
    
    # Build Hanzo binary first
    echo "Building Hanzo binary..."
    cd ..
    make
    cd s3-rdma-sidecar
    
    # Copy binary for Docker builds
    mkdir -p bin
    if [[ -f "../s3" ]]; then
        cp ../s3 bin/
    elif [[ -f "../bin/s3" ]]; then
        cp ../bin/s3 bin/
    elif [[ -f "../build/s3" ]]; then
        cp ../build/s3 bin/
    else
        echo "Error: Cannot find s3 binary"
        find .. -name "s3" -type f
        exit 1
    fi
    
    # Build RDMA sidecar
    echo "Building RDMA sidecar..."
    go build -o bin/demo-server cmd/sidecar/main.go
    
    # Build Docker images
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" build
    
    echo -e "${GREEN}✅ Images built successfully${NC}"
}

# Function to start services
start_services() {
    echo -e "${BLUE}🚀 Starting Hanzo RDMA Mount environment...${NC}"
    
    # Build images if needed
    if [[ ! -f "bin/s3" ]] || [[ ! -f "bin/demo-server" ]]; then
        build_images
    fi
    
    # Start services
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" up -d
    
    echo -e "${GREEN}✅ Services started${NC}"
    echo ""
    echo "Services are starting up. Use '$0 status' to check their status."
    echo "Use '$0 logs' to see the logs."
}

# Function to stop services
stop_services() {
    echo -e "${BLUE}🛑 Stopping Hanzo RDMA Mount environment...${NC}"
    
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down
    
    echo -e "${GREEN}✅ Services stopped${NC}"
}

# Function to restart services
restart_services() {
    echo -e "${BLUE}🔄 Restarting Hanzo RDMA Mount environment...${NC}"
    
    stop_services
    sleep 2
    start_services
}

# Function to show status
show_status() {
    echo -e "${BLUE}📊 Service Status${NC}"
    echo "================"
    
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps
    
    echo ""
    echo -e "${BLUE}🔍 Health Checks${NC}"
    echo "==============="
    
    # Check individual services
    check_service_health "Hanzo Master" "http://localhost:9333/cluster/status"
    check_service_health "Hanzo Volume" "http://localhost:8080/status"
    check_service_health "Hanzo Filer" "http://localhost:8888/"
    check_service_health "RDMA Sidecar" "http://localhost:8081/health"
    
    # Check mount status
    echo -n "Hanzo Mount: "
    if docker exec "${PROJECT_NAME}-hanzo-mount-1" mountpoint -q /mnt/hanzo 2>/dev/null; then
        echo -e "${GREEN}✅ Mounted${NC}"
    else
        echo -e "${RED}❌ Not mounted${NC}"
    fi
}

# Function to check service health
check_service_health() {
    local service_name=$1
    local health_url=$2
    
    echo -n "$service_name: "
    if curl -s "$health_url" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Healthy${NC}"
    else
        echo -e "${RED}❌ Unhealthy${NC}"
    fi
}

# Function to show logs
show_logs() {
    local service=$1
    
    if [[ -n "$service" ]]; then
        echo -e "${BLUE}📋 Logs for $service${NC}"
        echo "===================="
        $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" logs -f "$service"
    else
        echo -e "${BLUE}📋 Logs for all services${NC}"
        echo "======================="
        $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" logs -f
    fi
}

# Function to run integration tests
run_integration_tests() {
    echo -e "${BLUE}🧪 Running integration tests...${NC}"
    
    # Make sure services are running
    if ! $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps | grep -q "Up"; then
        echo -e "${RED}❌ Services are not running. Start them first with '$0 start'${NC}"
        exit 1
    fi
    
    # Run integration tests
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" --profile test run --rm integration-test
    
    # Show results
    if [[ -d "./test-results" ]]; then
        echo -e "${BLUE}📊 Test Results${NC}"
        echo "==============="
        
        if [[ -f "./test-results/overall.result" ]]; then
            local result
            result=$(cat "./test-results/overall.result")
            if [[ "$result" == "SUCCESS" ]]; then
                echo -e "${GREEN}🎉 ALL TESTS PASSED!${NC}"
            else
                echo -e "${RED}💥 SOME TESTS FAILED!${NC}"
            fi
        fi
        
        echo ""
        echo "Detailed results available in: ./test-results/"
        ls -la ./test-results/
    fi
}

# Function to run performance tests
run_performance_tests() {
    echo -e "${BLUE}🏁 Running performance tests...${NC}"
    
    # Make sure services are running
    if ! $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps | grep -q "Up"; then
        echo -e "${RED}❌ Services are not running. Start them first with '$0 start'${NC}"
        exit 1
    fi
    
    # Run performance tests
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" --profile performance run --rm performance-test
    
    # Show results
    if [[ -d "./performance-results" ]]; then
        echo -e "${BLUE}📊 Performance Results${NC}"
        echo "======================"
        echo ""
        echo "Results available in: ./performance-results/"
        ls -la ./performance-results/
        
        if [[ -f "./performance-results/performance_report.html" ]]; then
            echo ""
            echo -e "${GREEN}📄 HTML Report: ./performance-results/performance_report.html${NC}"
        fi
    fi
}

# Function to open shell in mount container
open_shell() {
    echo -e "${BLUE}🐚 Opening shell in mount container...${NC}"
    
    if ! $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps hanzo-mount | grep -q "Up"; then
        echo -e "${RED}❌ Mount container is not running${NC}"
        exit 1
    fi
    
    docker exec -it "${PROJECT_NAME}-hanzo-mount-1" /bin/bash
}

# Function to cleanup everything
cleanup_all() {
    echo -e "${BLUE}🧹 Full cleanup...${NC}"
    
    # Stop services
    $DOCKER_COMPOSE -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down -v --remove-orphans
    
    # Remove images
    echo "Removing Docker images..."
    docker images | grep "$PROJECT_NAME" | awk '{print $3}' | xargs -r docker rmi -f
    
    # Clean up local files
    rm -rf bin/ test-results/ performance-results/
    
    echo -e "${GREEN}✅ Full cleanup completed${NC}"
}

# Main function
main() {
    local command=${1:-""}
    
    # Check Docker Compose availability
    check_docker_compose
    
    case "$command" in
        "start")
            start_services
            ;;
        "stop")
            stop_services
            ;;
        "restart")
            restart_services
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs "${2:-}"
            ;;
        "test")
            run_integration_tests
            ;;
        "perf")
            run_performance_tests
            ;;
        "shell")
            open_shell
            ;;
        "cleanup")
            cleanup_all
            ;;
        "build")
            build_images
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        "")
            show_usage
            ;;
        *)
            echo -e "${RED}❌ Unknown command: $command${NC}"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"

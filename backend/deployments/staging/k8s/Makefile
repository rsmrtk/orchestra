# =============================================================================
# Makefile для Orchestra Backend Kubernetes Deployment (Staging)
# =============================================================================
# Використання: make <command>
# Приклади:
#   make deploy          # Deploy всього
#   make logs            # Дивитися логи
#   make status          # Статус всіх ресурсів
#   make help            # Показати всі команди
# =============================================================================

# Змінні
NAMESPACE := orchestra-staging
APP_NAME := orchestra-backend
DEPLOYMENT_NAME := orchestra-backend
SERVICE_NAME := orchestra-backend-service
INGRESS_NAME := orchestra-backend-ingress
HPA_NAME := orchestra-backend-hpa
PDB_NAME := orchestra-backend-pdb

# Docker образ (зміни на свій!)
IMAGE_REGISTRY := your-registry
IMAGE_NAME := orchestra-backend
IMAGE_TAG := staging

# Кольори для виводу
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: help
help: ## Показати цю допомогу
	@echo "$(GREEN)Orchestra Backend Kubernetes Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

# =============================================================================
# DEPLOY КОМАНДИ
# =============================================================================

.PHONY: deploy
deploy: ## Deploy всіх ресурсів (kustomize)
	@echo "$(GREEN)Deploying всіх ресурсів в $(NAMESPACE)...$(NC)"
	kubectl apply -k .
	@echo "$(GREEN)✓ Deploy завершено!$(NC)"

.PHONY: deploy-files
deploy-files: ## Deploy з окремих файлів
	@echo "$(GREEN)Creating namespace...$(NC)"
	kubectl apply -f namespace.yaml
	@echo "$(GREEN)Creating configmap...$(NC)"
	kubectl apply -f configmap.yaml
	@echo "$(GREEN)Creating secrets...$(NC)"
	kubectl apply -f secret.yaml
	@echo "$(GREEN)Creating deployment...$(NC)"
	kubectl apply -f deployment.yaml
	@echo "$(GREEN)Creating service...$(NC)"
	kubectl apply -f service.yaml
	@echo "$(GREEN)Creating ingress...$(NC)"
	kubectl apply -f ingress.yaml
	@echo "$(GREEN)Creating HPA...$(NC)"
	kubectl apply -f hpa.yaml
	@echo "$(GREEN)Creating PDB...$(NC)"
	kubectl apply -f pdb.yaml
	@echo "$(GREEN)Creating NetworkPolicy...$(NC)"
	kubectl apply -f network-policy.yaml
	@echo "$(GREEN)✓ Всі ресурси створені!$(NC)"

.PHONY: delete
delete: ## Видалити всі ресурси
	@echo "$(RED)Видалення всіх ресурсів з $(NAMESPACE)...$(NC)"
	@read -p "Ти впевнений? [y/N]: " confirm && [ "$$confirm" = "y" ] || exit 1
	kubectl delete -k .
	@echo "$(RED)✓ Всі ресурси видалені!$(NC)"

.PHONY: delete-namespace
delete-namespace: ## Видалити весь namespace (УВАГА: видалить ВСЕ!)
	@echo "$(RED)УВАГА: Це видалить ВСЕ в namespace $(NAMESPACE)!$(NC)"
	@read -p "Ти ДІЙСНО впевнений? [y/N]: " confirm && [ "$$confirm" = "y" ] || exit 1
	kubectl delete namespace $(NAMESPACE)
	@echo "$(RED)✓ Namespace видалено!$(NC)"

.PHONY: restart
restart: ## Перезапустити deployment
	@echo "$(YELLOW)Перезапуск $(DEPLOYMENT_NAME)...$(NC)"
	kubectl rollout restart deployment/$(DEPLOYMENT_NAME) -n $(NAMESPACE)
	@echo "$(GREEN)✓ Перезапуск ініційовано!$(NC)"

# =============================================================================
# STATUS КОМАНДИ
# =============================================================================

.PHONY: status
status: ## Статус всіх ресурсів
	@echo "$(GREEN)=== Статус $(NAMESPACE) ===$(NC)"
	kubectl get all -n $(NAMESPACE)

.PHONY: pods
pods: ## Список Pod'ів
	kubectl get pods -n $(NAMESPACE) -o wide

.PHONY: watch-pods
watch-pods: ## Watch Pod'ів (автооновлення)
	kubectl get pods -n $(NAMESPACE) --watch

.PHONY: services
services: ## Список Services
	kubectl get svc -n $(NAMESPACE)

.PHONY: ingress
ingress: ## Статус Ingress
	kubectl get ingress -n $(NAMESPACE)

.PHONY: hpa
hpa: ## Статус HPA
	kubectl get hpa -n $(NAMESPACE)

.PHONY: pdb
pdb: ## Статус PDB
	kubectl get pdb -n $(NAMESPACE)

.PHONY: endpoints
endpoints: ## Endpoints Service (IP Pod'ів)
	kubectl get endpoints $(SERVICE_NAME) -n $(NAMESPACE)

.PHONY: events
events: ## Події в namespace (останні)
	kubectl get events -n $(NAMESPACE) --sort-by='.lastTimestamp' | tail -20

# =============================================================================
# LOGS & DEBUG
# =============================================================================

.PHONY: logs
logs: ## Логи deployment (follow)
	kubectl logs -f deployment/$(DEPLOYMENT_NAME) -n $(NAMESPACE)

.PHONY: logs-tail
logs-tail: ## Останні 100 рядків логів
	kubectl logs deployment/$(DEPLOYMENT_NAME) -n $(NAMESPACE) --tail=100

.PHONY: logs-pod
logs-pod: ## Логи конкретного Pod'а (треба вказати POD_NAME=xxx)
	@if [ -z "$(POD_NAME)" ]; then \
		echo "$(RED)Помилка: вкажи POD_NAME=xxx$(NC)"; \
		echo "$(YELLOW)Доступні Pod'и:$(NC)"; \
		kubectl get pods -n $(NAMESPACE) -o name | sed 's/pod\///'; \
		exit 1; \
	fi
	kubectl logs -f $(POD_NAME) -n $(NAMESPACE)

.PHONY: describe-pod
describe-pod: ## Describe Pod (треба вказати POD_NAME=xxx)
	@if [ -z "$(POD_NAME)" ]; then \
		echo "$(RED)Помилка: вкажи POD_NAME=xxx$(NC)"; \
		exit 1; \
	fi
	kubectl describe pod $(POD_NAME) -n $(NAMESPACE)

.PHONY: exec
exec: ## Exec в Pod (shell) - треба вказати POD_NAME=xxx
	@if [ -z "$(POD_NAME)" ]; then \
		echo "$(RED)Помилка: вкажи POD_NAME=xxx$(NC)"; \
		exit 1; \
	fi
	kubectl exec -it $(POD_NAME) -n $(NAMESPACE) -- sh

.PHONY: port-forward
port-forward: ## Port forward Service (8080 -> 80)
	@echo "$(GREEN)Port forwarding $(SERVICE_NAME):80 -> localhost:8080$(NC)"
	@echo "$(YELLOW)Доступ: http://localhost:8080$(NC)"
	kubectl port-forward svc/$(SERVICE_NAME) 8080:80 -n $(NAMESPACE)

.PHONY: top
top: ## Використання ресурсів Pod'ами
	kubectl top pods -n $(NAMESPACE)

.PHONY: top-nodes
top-nodes: ## Використання ресурсів Node'ами
	kubectl top nodes

# =============================================================================
# SCALE & UPDATE
# =============================================================================

.PHONY: scale
scale: ## Масштабувати (треба вказати REPLICAS=N)
	@if [ -z "$(REPLICAS)" ]; then \
		echo "$(RED)Помилка: вкажи REPLICAS=N$(NC)"; \
		echo "$(YELLOW)Приклад: make scale REPLICAS=5$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Масштабування до $(REPLICAS) реплік...$(NC)"
	kubectl scale deployment/$(DEPLOYMENT_NAME) --replicas=$(REPLICAS) -n $(NAMESPACE)
	@echo "$(GREEN)✓ Масштабування завершено!$(NC)"

.PHONY: update-image
update-image: ## Оновити Docker образ (треба вказати TAG=xxx)
	@if [ -z "$(TAG)" ]; then \
		echo "$(RED)Помилка: вкажи TAG=xxx$(NC)"; \
		echo "$(YELLOW)Приклад: make update-image TAG=v1.2.3$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Оновлення образу до $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(TAG)...$(NC)"
	kubectl set image deployment/$(DEPLOYMENT_NAME) \
		backend=$(IMAGE_REGISTRY)/$(IMAGE_NAME):$(TAG) \
		-n $(NAMESPACE)
	@echo "$(GREEN)✓ Образ оновлено!$(NC)"
	@echo "$(YELLOW)Перевіряй статус: make rollout-status$(NC)"

.PHONY: rollout-status
rollout-status: ## Статус rollout
	kubectl rollout status deployment/$(DEPLOYMENT_NAME) -n $(NAMESPACE)

.PHONY: rollout-history
rollout-history: ## Історія rollout
	kubectl rollout history deployment/$(DEPLOYMENT_NAME) -n $(NAMESPACE)

.PHONY: rollout-undo
rollout-undo: ## Відкат до попередньої версії
	@echo "$(YELLOW)Відкат deployment...$(NC)"
	kubectl rollout undo deployment/$(DEPLOYMENT_NAME) -n $(NAMESPACE)
	@echo "$(GREEN)✓ Відкат завершено!$(NC)"

# =============================================================================
# CONFIG & SECRETS
# =============================================================================

.PHONY: update-config
update-config: ## Оновити ConfigMap
	@echo "$(YELLOW)Оновлення ConfigMap...$(NC)"
	kubectl apply -f configmap.yaml
	@echo "$(GREEN)✓ ConfigMap оновлено!$(NC)"
	@echo "$(YELLOW)Перезапуск Pod'ів для підхоплення змін...$(NC)"
	kubectl rollout restart deployment/$(DEPLOYMENT_NAME) -n $(NAMESPACE)

.PHONY: update-secret
update-secret: ## Оновити Secret
	@echo "$(YELLOW)Оновлення Secret...$(NC)"
	kubectl apply -f secret.yaml
	@echo "$(GREEN)✓ Secret оновлено!$(NC)"
	@echo "$(YELLOW)Перезапуск Pod'ів для підхоплення змін...$(NC)"
	kubectl rollout restart deployment/$(DEPLOYMENT_NAME) -n $(NAMESPACE)

.PHONY: get-config
get-config: ## Подивитися ConfigMap
	kubectl get configmap orchestra-backend-config -n $(NAMESPACE) -o yaml

.PHONY: get-secret
get-secret: ## Подивитися Secret (decoded)
	@echo "$(YELLOW)Secret (decoded):$(NC)"
	kubectl get secret orchestra-backend-secret -n $(NAMESPACE) -o json | \
		jq -r '.data | to_entries[] | "\(.key): \(.value | @base64d)"'

# =============================================================================
# BUILD & PUSH (опційно)
# =============================================================================

.PHONY: docker-build
docker-build: ## Build Docker образу
	@echo "$(GREEN)Building Docker образу...$(NC)"
	cd ../../.. && docker build -f deployments/staging/Dockerfile \
		-t $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "$(GREEN)✓ Образ побудовано: $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)$(NC)"

.PHONY: docker-push
docker-push: ## Push Docker образу
	@echo "$(GREEN)Pushing Docker образу...$(NC)"
	docker push $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	@echo "$(GREEN)✓ Образ запушено!$(NC)"

.PHONY: docker-build-push
docker-build-push: docker-build docker-push ## Build і Push образу

# =============================================================================
# TESTING
# =============================================================================

.PHONY: test-api
test-api: ## Тестувати API через port-forward
	@echo "$(YELLOW)Тестування API через port-forward...$(NC)"
	@(kubectl port-forward svc/$(SERVICE_NAME) 8080:80 -n $(NAMESPACE) > /dev/null 2>&1 &); \
	sleep 2; \
	curl -s http://localhost:8080/congratulation | jq .; \
	pkill -f "kubectl port-forward"

.PHONY: test-ingress
test-ingress: ## Тестувати через Ingress (треба вказати HOST=xxx)
	@if [ -z "$(HOST)" ]; then \
		echo "$(RED)Помилка: вкажи HOST=xxx$(NC)"; \
		echo "$(YELLOW)Приклад: make test-ingress HOST=staging-api.orchestra.example.com$(NC)"; \
		exit 1; \
	fi
	curl -s https://$(HOST)/api/congratulation | jq .

.PHONY: load-test
load-test: ## Load test (Apache Bench) - треба вказати URL=xxx
	@if [ -z "$(URL)" ]; then \
		echo "$(RED)Помилка: вкажи URL=xxx$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)Load testing $(URL)...$(NC)"
	ab -n 10000 -c 100 $(URL)

# =============================================================================
# CLEANUP
# =============================================================================

.PHONY: clean-failed
clean-failed: ## Видалити failed Pod'и
	@echo "$(YELLOW)Видалення failed Pod'ів...$(NC)"
	kubectl delete pods -n $(NAMESPACE) --field-selector=status.phase=Failed
	@echo "$(GREEN)✓ Failed Pod'и видалені!$(NC)"

.PHONY: clean-completed
clean-completed: ## Видалити completed Pod'и
	@echo "$(YELLOW)Видалення completed Pod'ів...$(NC)"
	kubectl delete pods -n $(NAMESPACE) --field-selector=status.phase=Succeeded
	@echo "$(GREEN)✓ Completed Pod'и видалені!$(NC)"

# =============================================================================
# UTILITIES
# =============================================================================

.PHONY: diff
diff: ## Diff (що зміниться)
	kubectl diff -k .

.PHONY: dry-run
dry-run: ## Dry run (без застосування)
	kubectl apply -k . --dry-run=client -o yaml

.PHONY: context
context: ## Показати поточний context
	@echo "$(YELLOW)Поточний context:$(NC)"
	kubectl config current-context
	@echo ""
	@echo "$(YELLOW)Поточний namespace:$(NC)"
	kubectl config view --minify --output 'jsonpath={..namespace}'
	@echo ""

.PHONY: switch-context
switch-context: ## Перемкнути context (треба вказати CONTEXT=xxx)
	@if [ -z "$(CONTEXT)" ]; then \
		echo "$(RED)Помилка: вкажи CONTEXT=xxx$(NC)"; \
		echo "$(YELLOW)Доступні contexts:$(NC)"; \
		kubectl config get-contexts; \
		exit 1; \
	fi
	kubectl config use-context $(CONTEXT)

.PHONY: yaml
yaml: ## Згенерувати фінальний YAML (kustomize build)
	kubectl kustomize .

# =============================================================================
# MONITORING
# =============================================================================

.PHONY: watch-hpa
watch-hpa: ## Watch HPA (автомасштабування)
	kubectl get hpa $(HPA_NAME) -n $(NAMESPACE) --watch

.PHONY: watch-all
watch-all: ## Watch всіх ресурсів
	watch -n 2 'kubectl get all,ingress,hpa,pdb -n $(NAMESPACE)'

# =============================================================================
# QUICK ACTIONS
# =============================================================================

.PHONY: quick-deploy
quick-deploy: docker-build-push deploy ## Build, Push, Deploy (все разом)
	@echo "$(GREEN)✓ Швидкий deploy завершено!$(NC)"

.PHONY: quick-update
quick-update: ## Швидке оновлення (build, push, update image)
	@$(MAKE) docker-build TAG=staging-$(shell date +%s)
	@$(MAKE) docker-push TAG=staging-$(shell date +%s)
	@$(MAKE) update-image TAG=staging-$(shell date +%s)

# =============================================================================
# INFO
# =============================================================================

.PHONY: info
info: ## Показати інформацію про deployment
	@echo "$(GREEN)╔════════════════════════════════════════════════════╗$(NC)"
	@echo "$(GREEN)║  Orchestra Backend - Staging Environment          ║$(NC)"
	@echo "$(GREEN)╚════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@echo "$(YELLOW)Namespace:$(NC)       $(NAMESPACE)"
	@echo "$(YELLOW)App Name:$(NC)        $(APP_NAME)"
	@echo "$(YELLOW)Deployment:$(NC)      $(DEPLOYMENT_NAME)"
	@echo "$(YELLOW)Service:$(NC)         $(SERVICE_NAME)"
	@echo "$(YELLOW)Ingress:$(NC)         $(INGRESS_NAME)"
	@echo "$(YELLOW)HPA:$(NC)             $(HPA_NAME)"
	@echo "$(YELLOW)PDB:$(NC)             $(PDB_NAME)"
	@echo ""
	@echo "$(YELLOW)Docker Image:$(NC)    $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)"
	@echo ""
	@echo "$(YELLOW)Використання:$(NC)    make help"
	@echo ""

# Default command
.DEFAULT_GOAL := help

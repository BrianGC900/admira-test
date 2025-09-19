
### SYSTEM_DESIGN.md

```markdown
# System Design - Admira Backend

## Idempotencia & Reprocesamiento

**Idempotencia**: El endpoint `/ingest/run` está diseñado para ser idempotente mediante el parámetro `since`. Al especificar una fecha, se reprocesarán solo los datos desde esa fecha, evitando duplicados.

**Reprocesamiento**: Se puede reprocesar datos históricos cambiando el parámetro `since`. En producción, se implementaría un sistema de checkpointing para tracking del último procesamiento.

## Particionamiento & Retención

**Particionamiento**: Los datos se particionan naturalmente por fecha, facilitando consultas por rangos temporales.

**Retención**: En la implementación actual con almacenamiento en memoria, los datos se pierden al reiniciar. En producción, se usaría una base de datos con políticas de retención configuradas (ej: 90 días para datos detallados, agregados mensuales permanentes).

## Concurrencia & Throughput

**Concurrencia**: 
- Uso de goroutines para procesamiento paralelo de diferentes UTMs
- Worker pools para procesamiento de grandes volúmenes de datos

**Throughput**: 
- El procesamiento actual es sincrónico por simplicidad
- En producción, se implementaría un sistema asíncrono con colas (RabbitMQ, Kafka) para alta escalabilidad

## Calidad de datos

**UTMs ausentes**: 
- Campos UTM faltantes se normalizan a "unknown"
- Estrategia de fallback: Agrupar por combinación disponible de UTM parameters

**Validaciones**:
- Validación de formatos de fecha
- Verificación de campos requeridos
- Sanitización de valores nulos o incorrectos

## Observabilidad

**Logging**:
- Logs estructurados en JSON
- Campos: timestamp, level, method, path, status, duration, correlation_id

**Métricas** (opcional):
- Contadores: requests procesados, errores
- Latencia: tiempos de procesamiento por etapa ETL
- Métricas de negocio: ROAS, CPA por canal

**Trazabilidad**:
- Correlation ID por request para tracking entre servicios

## Evolución en el ecosistema Admira

**Data Lake/ETL**:
- Este servicio puede ser el primer eslabón de un pipeline ETL más complejo
- Exportación de datos procesados a un data lake (S3, BigQuery)
- Integración con herramientas de BI (Tableau, Looker)

**Contratos de API**:
- Versionado de APIs para compatibilidad hacia atrás
- Documentación OpenAPI/Swagger
- Schemas Avro/Protobuf para serialización eficiente

**Escalabilidad**:
- Desacoplamiento con colas de mensajería
- Autoescalado basado en métricas de uso
- Cache distribuido (Redis) para consultas frecuentes
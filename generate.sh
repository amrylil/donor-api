#!/bin/bash

# Hentikan script jika ada error
set -e

# 1. Periksa apakah nama fitur diberikan
if [ -z "$1" ]; then
  echo "❌ Error: Nama fitur wajib diisi."
  echo "   Contoh: ./generate_crud.sh Product"
  echo "   Contoh: ./generate_crud.sh BloodRequest"
  exit 1
fi

# 2. Siapkan variabel nama dari input (e.g., BloodRequest)
FEATURE_NAME_INPUT=$1

# Membuat versi PascalCase (e.g., BloodRequest)
# Mengkapitalkan huruf pertama dari input, asumsi input adalah camelCase atau PascalCase.
FEATURE_NAME_PASCAL="$(tr '[:lower:]' '[:upper:]' <<< ${FEATURE_NAME_INPUT:0:1})${FEATURE_NAME_INPUT:1}"

# Membuat versi camelCase untuk nama variabel Go (e.g., bloodRequest)
# Mengubah huruf pertama dari PascalCase menjadi huruf kecil.
FEATURE_NAME_CAMEL="$(tr '[:upper:]' '[:lower:]' <<< ${FEATURE_NAME_PASCAL:0:1})${FEATURE_NAME_PASCAL:1}"

# Membuat versi snake_case untuk nama file (e.g., blood_request)
# Mengubah PascalCase menjadi snake_case.
FEATURE_NAME_SNAKE=$(echo "$FEATURE_NAME_PASCAL" | sed -e 's/\([A-Z]\)/_\l\1/g' -e 's/^_//')

# Membuat versi plural snake_case untuk URL (e.g., blood_requests)
FEATURE_NAME_SNAKE_PLURAL="${FEATURE_NAME_SNAKE}s"

# Membuat versi plural camelCase untuk nama variabel di Go (e.g., bloodRequests)
FEATURE_NAME_CAMEL_PLURAL="${FEATURE_NAME_CAMEL}s"


# Path direktori
ENTITY_PATH="internal/entity"
DTO_PATH="internal/delivery/http/dto"
HELPER_PATH="internal/delivery/http/helper"
REPO_PATH="internal/repository"
PERSISTENCE_PATH="internal/infrastructure/persistence"
USECASE_PATH="internal/usecase"
HANDLER_PATH="internal/delivery/http/handler"
ROUTES_PATH="internal/delivery/routes"

echo "🚀 Mulai membuat fitur CRUD mandiri untuk: $FEATURE_NAME_PASCAL"

# Membuat direktori jika belum ada
mkdir -p $ENTITY_PATH
mkdir -p $DTO_PATH
mkdir -p $HELPER_PATH
mkdir -p $REPO_PATH
mkdir -p $PERSISTENCE_PATH
mkdir -p $USECASE_PATH
mkdir -p $HANDLER_PATH
mkdir -p $ROUTES_PATH


# 3. Buat file Entity
cat <<EOF > ./${ENTITY_PATH}/${FEATURE_NAME_SNAKE}.go
package entity

import (
  "time"

  "github.com/google/uuid"
  "gorm.io/gorm"
)

type ${FEATURE_NAME_PASCAL} struct {
  ID        uuid.UUID      \`gorm:"type:uuid;primary_key;" json:"id"\`
  Title     string         \`gorm:"type:varchar(255)" json:"title"\`
  CreatedAt time.Time      \`json:"created_at"\`
  UpdatedAt time.Time      \`json:"updated_at"\`
  DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\`
}

func (p *${FEATURE_NAME_PASCAL}) BeforeCreate(tx *gorm.DB) (err error) {
  p.ID = uuid.New()
  return
}
EOF
echo "✅ Dibuat: ${ENTITY_PATH}/${FEATURE_NAME_SNAKE}.go"

# 4. Buat file DTO
cat <<EOF > ./${DTO_PATH}/${FEATURE_NAME_SNAKE}_dto.go
package dto

import "time"

// DTO untuk request body (Create & Update)
type ${FEATURE_NAME_PASCAL}Request struct {
  Title string \`json:"title" binding:"required"\`
}

// DTO untuk response (data aman untuk publik)
type ${FEATURE_NAME_PASCAL}Response struct {
  ID        string    \`json:"id"\`
  Title     string    \`json:"title"\`
  CreatedAt time.Time \`json:"created_at"\`
  UpdatedAt time.Time \`json:"updated_at"\`
}
EOF
echo "✅ Dibuat: ${DTO_PATH}/${FEATURE_NAME_SNAKE}_dto.go"


# 5. Buat file DTO & Helper Generik (jika belum ada)
RESPONSE_DTO_FILE="./${DTO_PATH}/response_dto.go"
if [ ! -f "$RESPONSE_DTO_FILE" ]; then
    cat <<EOF > $RESPONSE_DTO_FILE
package dto

// APIResponse adalah template dasar untuk semua respons JSON dari API.
type APIResponse[T any] struct {
  Success bool   \`json:"success"\`
  Message string \`json:"message"\`
  Data    T      \`json:"data,omitempty"\`
  Error   any    \`json:"error,omitempty"\`
}
EOF
    echo "✅ Dibuat: ${DTO_PATH}/response_dto.go"
else
    echo "☑️  Dilewati: ${DTO_PATH}/response_dto.go sudah ada"
fi

PAGINATION_FILE="./${DTO_PATH}/pagination_dto.go"
if [ ! -f "$PAGINATION_FILE" ]; then
    cat <<EOF > $PAGINATION_FILE
package dto

// PaginatedResponse adalah struktur DTO generik untuk respons pagination.
type PaginatedResponse[T any] struct {
  Data       []T   \`json:"data"\`
  TotalItems int64 \`json:"total_items"\`
  Page       int   \`json:"page"\`
  Limit      int   \`json:"limit"\`
}
EOF
    echo "✅ Dibuat: ${DTO_PATH}/pagination_dto.go"
else
    echo "☑️  Dilewati: ${DTO_PATH}/pagination_dto.go sudah ada"
fi

HELPER_FILE="./${HELPER_PATH}/response_helper.go"
if [ ! -f "$HELPER_FILE" ]; then
    cat <<EOF > $HELPER_FILE
package helper

import (
  "donor-api/internal/delivery/http/dto"
  "github.com/gin-gonic/gin"
)

func SendSuccessResponse[T any](c *gin.Context, statusCode int, message string, data T) {
  response := dto.APIResponse[T]{
    Success: true,
    Message: message,
    Data:    data,
  }
  c.JSON(statusCode, response)
}

func SendErrorResponse(c *gin.Context, statusCode int, message string) {
    response := dto.APIResponse[any]{
        Success: false,
        Message: message,
    }
    c.JSON(statusCode, response)
}
EOF
    echo "✅ Dibuat: ${HELPER_PATH}/response_helper.go"
else
    echo "☑️  Dilewati: ${HELPER_PATH}/response_helper.go sudah ada"
fi


# 6. Buat Repository (Interface)
cat <<EOF > ./${REPO_PATH}/${FEATURE_NAME_SNAKE}_repository.go
package repository

import (
  "context"
  "donor-api/internal/entity"

  "github.com/google/uuid"
)

type ${FEATURE_NAME_PASCAL}Repository interface {
  Save(ctx context.Context, ${FEATURE_NAME_CAMEL} *entity.${FEATURE_NAME_PASCAL}) error
  FindAll(ctx context.Context, limit, offset int) ([]entity.${FEATURE_NAME_PASCAL}, int64, error)
  FindByID(ctx context.Context, id uuid.UUID) (entity.${FEATURE_NAME_PASCAL}, error)
  Update(ctx context.Context, ${FEATURE_NAME_CAMEL} entity.${FEATURE_NAME_PASCAL}) (entity.${FEATURE_NAME_PASCAL}, error)
  Delete(ctx context.Context, id uuid.UUID) error
}
EOF
echo "✅ Dibuat: ${REPO_PATH}/${FEATURE_NAME_SNAKE}_repository.go"


# 7. Buat Repository (Implementation)
cat <<EOF > ./${PERSISTENCE_PATH}/${FEATURE_NAME_SNAKE}_repository_impl.go
package persistence

import (
  "context"
  "donor-api/internal/entity"
  "donor-api/internal/repository"

  "github.com/google/uuid"
  "gorm.io/gorm"
)

type ${FEATURE_NAME_CAMEL}RepositoryImpl struct {
  db *gorm.DB
}

func New${FEATURE_NAME_PASCAL}Repository(db *gorm.DB) repository.${FEATURE_NAME_PASCAL}Repository {
  return &${FEATURE_NAME_CAMEL}RepositoryImpl{db: db}
}

func (r *${FEATURE_NAME_CAMEL}RepositoryImpl) Save(ctx context.Context, ${FEATURE_NAME_CAMEL} *entity.${FEATURE_NAME_PASCAL}) error {
  return r.db.WithContext(ctx).Create(${FEATURE_NAME_CAMEL}).Error
}

func (r *${FEATURE_NAME_CAMEL}RepositoryImpl) FindAll(ctx context.Context, limit, offset int) ([]entity.${FEATURE_NAME_PASCAL}, int64, error) {
  var ${FEATURE_NAME_CAMEL_PLURAL} []entity.${FEATURE_NAME_PASCAL}
  var total int64

  if err := r.db.WithContext(ctx).Model(&entity.${FEATURE_NAME_PASCAL}{}).Count(&total).Error; err != nil {
    return nil, 0, err
  }

  if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&${FEATURE_NAME_CAMEL_PLURAL}).Error; err != nil {
    return nil, 0, err
  }

  return ${FEATURE_NAME_CAMEL_PLURAL}, total, nil
}

func (r *${FEATURE_NAME_CAMEL}RepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (entity.${FEATURE_NAME_PASCAL}, error) {
  var ${FEATURE_NAME_CAMEL} entity.${FEATURE_NAME_PASCAL}
  // GORM dapat mencari berdasarkan primary key secara langsung.
  err := r.db.WithContext(ctx).First(&${FEATURE_NAME_CAMEL}, id).Error
  return ${FEATURE_NAME_CAMEL}, err
}

func (r *${FEATURE_NAME_CAMEL}RepositoryImpl) Update(ctx context.Context, ${FEATURE_NAME_CAMEL} entity.${FEATURE_NAME_PASCAL}) (entity.${FEATURE_NAME_PASCAL}, error) {
  err := r.db.WithContext(ctx).Save(&${FEATURE_NAME_CAMEL}).Error
  return ${FEATURE_NAME_CAMEL}, err
}

func (r *${FEATURE_NAME_CAMEL}RepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
  // GORM dapat menghapus berdasarkan primary key secara langsung.
  return r.db.WithContext(ctx).Delete(&entity.${FEATURE_NAME_PASCAL}{}, id).Error
}
EOF
echo "✅ Dibuat: ${PERSISTENCE_PATH}/${FEATURE_NAME_SNAKE}_repository_impl.go"


# 8. Buat Usecase (Menggunakan Copier)
cat <<EOF > ./${USECASE_PATH}/${FEATURE_NAME_SNAKE}_usecase.go
package usecase

import (
  "context"
  "donor-api/internal/delivery/http/dto"
  "donor-api/internal/entity"
  "donor-api/internal/repository"

  "github.com/google/uuid"
  "github.com/jinzhu/copier"
)

// --- Interface ---
type ${FEATURE_NAME_PASCAL}Usecase interface {
  Create(ctx context.Context, req dto.${FEATURE_NAME_PASCAL}Request) (dto.${FEATURE_NAME_PASCAL}Response, error)
  FindAll(ctx context.Context, page, limit int) (dto.PaginatedResponse[dto.${FEATURE_NAME_PASCAL}Response], error)
  FindByID(ctx context.Context, id uuid.UUID) (dto.${FEATURE_NAME_PASCAL}Response, error)
  Update(ctx context.Context, id uuid.UUID, req dto.${FEATURE_NAME_PASCAL}Request) (dto.${FEATURE_NAME_PASCAL}Response, error)
  Delete(ctx context.Context, id uuid.UUID) error
}

// --- Implementation ---
type ${FEATURE_NAME_CAMEL}UsecaseImpl struct {
  repo repository.${FEATURE_NAME_PASCAL}Repository
}

func New${FEATURE_NAME_PASCAL}Usecase(repo repository.${FEATURE_NAME_PASCAL}Repository) ${FEATURE_NAME_PASCAL}Usecase {
  return &${FEATURE_NAME_CAMEL}UsecaseImpl{repo: repo}
}

func (uc *${FEATURE_NAME_CAMEL}UsecaseImpl) Create(ctx context.Context, req dto.${FEATURE_NAME_PASCAL}Request) (dto.${FEATURE_NAME_PASCAL}Response, error) {
  var ${FEATURE_NAME_CAMEL} entity.${FEATURE_NAME_PASCAL}
  var res dto.${FEATURE_NAME_PASCAL}Response

  copier.Copy(&${FEATURE_NAME_CAMEL}, &req)

  if err := uc.repo.Save(ctx, &${FEATURE_NAME_CAMEL}); err != nil {
    return res, err
  }

  // Salin field yang cocok, lalu atur ID secara manual.
  copier.Copy(&res, &${FEATURE_NAME_CAMEL})
  res.ID = ${FEATURE_NAME_CAMEL}.ID.String()
  return res, nil
}

func (uc *${FEATURE_NAME_CAMEL}UsecaseImpl) FindAll(ctx context.Context, page, limit int) (dto.PaginatedResponse[dto.${FEATURE_NAME_PASCAL}Response], error) {
  offset := (page - 1) * limit
  var paginatedResponse dto.PaginatedResponse[dto.${FEATURE_NAME_PASCAL}Response]

  items, total, err := uc.repo.FindAll(ctx, limit, offset)
  if err != nil {
    return paginatedResponse, err
  }

  var itemResponses []dto.${FEATURE_NAME_PASCAL}Response
  copier.Copy(&itemResponses, &items)

  // ID perlu di-mapping manual karena tipe berbeda (uuid.UUID -> string)
  for i := range items {
    itemResponses[i].ID = items[i].ID.String()
  }

  paginatedResponse = dto.PaginatedResponse[dto.${FEATURE_NAME_PASCAL}Response]{
    Data:       itemResponses,
    TotalItems: total,
    Page:       page,
    Limit:      limit,
  }
  return paginatedResponse, nil
}

func (uc *${FEATURE_NAME_CAMEL}UsecaseImpl) FindByID(ctx context.Context, id uuid.UUID) (dto.${FEATURE_NAME_PASCAL}Response, error) {
  var res dto.${FEATURE_NAME_PASCAL}Response
  ${FEATURE_NAME_CAMEL}, err := uc.repo.FindByID(ctx, id)
  if err != nil {
    return res, err
  }

  copier.Copy(&res, &${FEATURE_NAME_CAMEL})
  res.ID = ${FEATURE_NAME_CAMEL}.ID.String()
  return res, nil
}

func (uc *${FEATURE_NAME_CAMEL}UsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.${FEATURE_NAME_PASCAL}Request) (dto.${FEATURE_NAME_PASCAL}Response, error) {
  var res dto.${FEATURE_NAME_PASCAL}Response
  ${FEATURE_NAME_CAMEL}, err := uc.repo.FindByID(ctx, id)
  if err != nil {
    return res, err
  }

  copier.Copy(&${FEATURE_NAME_CAMEL}, &req)

  updated${FEATURE_NAME_PASCAL}, err := uc.repo.Update(ctx, ${FEATURE_NAME_CAMEL})
  if err != nil {
    return res, err
  }

  copier.Copy(&res, &updated${FEATURE_NAME_PASCAL})
  res.ID = updated${FEATURE_NAME_PASCAL}.ID.String()
  return res, nil
}

func (uc *${FEATURE_NAME_CAMEL}UsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
  _, err := uc.repo.FindByID(ctx, id)
  if err != nil {
    return err
  }
  return uc.repo.Delete(ctx, id)
}
EOF
echo "✅ Dibuat: ${USECASE_PATH}/${FEATURE_NAME_SNAKE}_usecase.go"


# 9. Buat Handler
cat <<EOF > ./${HANDLER_PATH}/${FEATURE_NAME_SNAKE}_handler.go
package handler

import (
  "donor-api/internal/delivery/http/dto"
  "donor-api/internal/delivery/http/helper"
  "donor-api/internal/usecase"
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
  "github.com/google/uuid"
)

type ${FEATURE_NAME_PASCAL}Handler struct {
  usecase usecase.${FEATURE_NAME_PASCAL}Usecase
}

func New${FEATURE_NAME_PASCAL}Handler(usecase usecase.${FEATURE_NAME_PASCAL}Usecase) *${FEATURE_NAME_PASCAL}Handler {
  return &${FEATURE_NAME_PASCAL}Handler{usecase: usecase}
}

func (h *${FEATURE_NAME_PASCAL}Handler) Create(c *gin.Context) {
  var req dto.${FEATURE_NAME_PASCAL}Request
  if err := c.ShouldBindJSON(&req); err != nil {
    helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
    return
  }

  res, err := h.usecase.Create(c.Request.Context(), req)
  if err != nil {
    helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
    return
  }

  helper.SendSuccessResponse(c, http.StatusCreated, "${FEATURE_NAME_PASCAL} created successfully", res)
}

func (h *${FEATURE_NAME_PASCAL}Handler) GetAll(c *gin.Context) {
  page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
  limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

  paginatedResponse, err := h.usecase.FindAll(c.Request.Context(), page, limit)
  if err != nil {
    helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
    return
  }

  helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved ${FEATURE_NAME_SNAKE_PLURAL}", paginatedResponse)
}

func (h *${FEATURE_NAME_PASCAL}Handler) GetByID(c *gin.Context) {
  id, err := uuid.Parse(c.Param("id"))
  if err != nil {
    helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
    return
  }
  res, err := h.usecase.FindByID(c.Request.Context(), id)
  if err != nil {
    helper.SendErrorResponse(c, http.StatusNotFound, "Record not found")
    return
  }

  helper.SendSuccessResponse(c, http.StatusOK, "Successfully retrieved ${FEATURE_NAME_SNAKE}", res)
}

func (h *${FEATURE_NAME_PASCAL}Handler) Update(c *gin.Context) {
  id, err := uuid.Parse(c.Param("id"))
  if err != nil {
    helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
    return
  }
  var req dto.${FEATURE_NAME_PASCAL}Request
  if err := c.ShouldBindJSON(&req); err != nil {
    helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
    return
  }

  res, err := h.usecase.Update(c.Request.Context(), id, req)
  if err != nil {
    helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
    return
  }

  helper.SendSuccessResponse(c, http.StatusOK, "${FEATURE_NAME_PASCAL} updated successfully", res)
}

func (h *${FEATURE_NAME_PASCAL}Handler) Delete(c *gin.Context) {
  id, err := uuid.Parse(c.Param("id"))
  if err != nil {
    helper.SendErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
    return
  }

  err = h.usecase.Delete(c.Request.Context(), id)
  if err != nil {
    helper.SendErrorResponse(c, http.StatusInternalServerError, err.Error())
    return
  }
  helper.SendSuccessResponse(c, http.StatusOK, "${FEATURE_NAME_PASCAL} deleted successfully", "")
}
EOF
echo "✅ Dibuat: ${HANDLER_PATH}/${FEATURE_NAME_SNAKE}_handler.go"


# 10. Buat Routes
cat <<EOF > ./${ROUTES_PATH}/${FEATURE_NAME_SNAKE}_routes.go
package routes

import (
  "donor-api/internal/delivery/http/handler"
  "github.com/gin-gonic/gin"
)

func Init${FEATURE_NAME_PASCAL}Routes(
  router *gin.RouterGroup,
  handler *handler.${FEATURE_NAME_PASCAL}Handler,
) {
  ${FEATURE_NAME_SNAKE_PLURAL}Routes := router.Group("/${FEATURE_NAME_SNAKE_PLURAL}")
  {
    ${FEATURE_NAME_SNAKE_PLURAL}Routes.POST("", handler.Create)
    ${FEATURE_NAME_SNAKE_PLURAL}Routes.GET("", handler.GetAll)
    ${FEATURE_NAME_SNAKE_PLURAL}Routes.GET("/:id", handler.GetByID)
    ${FEATURE_NAME_SNAKE_PLURAL}Routes.PUT("/:id", handler.Update)
    ${FEATURE_NAME_SNAKE_PLURAL}Routes.DELETE("/:id", handler.Delete)
  }
}
EOF
echo "✅ Dibuat: ${ROUTES_PATH}/${FEATURE_NAME_SNAKE}_routes.go"


package k8s

import (
	"github.com/go-playground/validator/v10"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/logger"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
)

type MetaDataDto struct {
	Name      string `json:"name" validate:"regexp,lt=40"`
	Namespace string `json:"namespace,omitempty"`
}

func (metaDto *MetaDataDto) ObjectMetaFromMetadataDto() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      metaDto.Name,
		Namespace: metaDto.Namespace,
	}
}

func (dto *MetaDataDto) Validate() *restErrors.RestErr {
	newValidator := validator.New()
	err := newValidator.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile("^([a-z]|[0-9])+([a-z]|[0-9]|-)+$")
		return re.MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Panic("USER_DTO_VALIDATE", err)
		return restErrors.NewInternalServerError("something went wrong!")
	}

	err = newValidator.Struct(dto)

	if err != nil {
		fields := map[string]string{}
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name":
				fields["name"] = "name should only contains a lowercase alphanumeric characters or special character (-) and maximum length of 40 chars"
				break
			}
		}

		if len(fields) > 0 {
			return restErrors.NewValidationError(fields)
		}
	}

	return nil
}

func DefaultResources(res *sharedAPI.Resources) {
	res.CPU = "1"
	res.Memory = "1Gi"
}

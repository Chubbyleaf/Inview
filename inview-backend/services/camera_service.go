package services

import (
	"context"
	"errors"
	"insense-local/data"
	"insense-local/models"
	"time"
)

type CameraServiceInterface interface {
	CreateCamera(c context.Context, camera *models.Camera) error
	FetchCameras(c context.Context) ([]models.Camera, error)
	DeleteCamera(c context.Context, id int, ts TaskServiceInterface) error
	UpdateCamera(c context.Context, id int, name, remark string) error
	FindCamera(c context.Context, id int) (models.Camera, error)
}

func CameraService(cameraData data.CameraDataInterface, timeout time.Duration) CameraServiceInterface {
	return &cameraService{
		cameraData:     cameraData,
		contextTimeout: timeout,
	}
}

type cameraService struct {
	cameraData     data.CameraDataInterface
	contextTimeout time.Duration
}

func (cs *cameraService) CreateCamera(c context.Context, camera *models.Camera) error {
	ctx, cancel := context.WithTimeout(c, cs.contextTimeout)
	defer cancel()

	// 获取最小的可用设备 ID
	deviceID, err := cs.cameraData.GetNextDeviceID(ctx)
	if err != nil {
		return err
	}

	camera.DeviceID = deviceID
	return cs.cameraData.Create(ctx, camera)
}

func (cs *cameraService) FetchCameras(c context.Context) ([]models.Camera, error) {
	ctx, cancel := context.WithTimeout(c, cs.contextTimeout)
	defer cancel()
	return cs.cameraData.Fetch(ctx)
}

func (cs *cameraService) DeleteCamera(c context.Context, id int, ts TaskServiceInterface) error {
	ctx, cancel := context.WithTimeout(c, cs.contextTimeout)
	defer cancel()
	exists, err := ts.DeviceExistTask(c, id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("该设备存在任务，无法删除！")
	}
	return cs.cameraData.Delete(ctx, id)
}

func (cs *cameraService) UpdateCamera(c context.Context, id int, name, remark string) error {
	ctx, cancel := context.WithTimeout(c, cs.contextTimeout)
	defer cancel()

	camera := &models.Camera{
		Name:   name,
		Remark: remark,
	}

	return cs.cameraData.Update(ctx, id, camera)
}

func (cs *cameraService) FindCamera(c context.Context, id int) (models.Camera, error) {
	ctx, cancel := context.WithTimeout(c, cs.contextTimeout)
	defer cancel()

	return cs.cameraData.FindByDeviceId(ctx, id)
}

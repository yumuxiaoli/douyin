package controller

import (
	"douyin/models"
	"douyin/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]models.User{
	"zhangleidouyin": {
		Model: gorm.Model{
			ID: 1,
		},
		Name:     "zhanglei",
		Password: "douyin",
		// Id:            1,
		// Name:          "zhanglei",
		// FollowCount:   10,
		// FollowerCount: 5,
		// IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User models.UserInfo `json:"user"`
}

func Register(c *fiber.Ctx) error {
	username := c.Query("username")
	password := c.Query("password")

	if _, err := service.GetUserByName(username); err == nil {
		fmt.Println("The suer exits")
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User already exist",
			},
		})
	}

	newUser := models.User{
		Name:     username,
		Password: password,
	}

	if err := service.CreateUser(&newUser); err != nil {
		fmt.Println("插入失败", err)
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 2,
				StatusMsg:  "User insertion error",
			},
		})
	}

	fmt.Println("插入成功")
	if token, err := service.GenerateToken(&newUser); err == nil {
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 0,
			},
			UserId: int64(newUser.ID),
			Token:  token,
		})
	}

	fmt.Println("创建token失败")
	return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
		Response: Response{
			StatusCode: 3,
			StatusMsg:  "Unable to create token",
		},
		UserId: int64(newUser.ID),
	})
}

func Login(c *fiber.Ctx) error {
	username := c.Query("username")
	password := c.Query("password")

	user, err := service.GetUserByName(username)

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
		})
	}

	if user.Password != password {
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 2,
				StatusMsg:  "Password doesn't match",
			},
		})
	}

	if token, err := service.GenerateToken(&user); err == nil {
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "Login successfully",
			},
			UserId: int64(user.ID),
			Token:  token,
		})
	}

	return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
		Response: Response{
			StatusCode: 2,
			StatusMsg:  "Unable to create token",
		},
	})
}

func UserInfo(c *fiber.Ctx) error {
	token := c.Query("token")
	uid, _ := strconv.Atoi(c.Query("user_id"))
	if _, err := service.ParseToken(token); err != nil {
		return c.Status(http.StatusOK).JSON(
			UserResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "user unauthorized",
				},
				
			},
		)
	}
	if user, err := service.GetUserById(uint(uid)); err != nil {
		return c.Status(fiber.StatusOK).JSON(
			UserResponse{
				Response: Response{StatusCode: 0},
				User: service.GenerateUserInfo(&user),
			},
		)
	}
	return c.Status(fiber.StatusOK).JSON(
		UserResponse{
			Response: Response{
				StatusCode: 2,
				StatusMsg:  "user not exist",
			},
		},
	)
}

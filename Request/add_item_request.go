package Request

type AddItem = struct {
    Name string `validate:"required"`
}
package outbox

import "testing"

func TestBuildOutboxTopicNamefromEventType(t *testing.T) {
	type args struct {
		eventType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "When RestaurantCreated then outbox-restaurant-created",
			args: args{eventType: "RestaurantCreated"},
			want: "outbox-restaurant-created",
		},
		{
			name: "When RestaurantDeleted then outbox-restaurant-deleted",
			args: args{eventType: "RestaurantDeleted"},
			want: "outbox-restaurant-deleted",
		},
		{
			name: "When RestaurantMenuUpdated then outbox-restaurant-menu-updated",
			args: args{eventType: "RestaurantMenuUpdated"},
			want: "outbox-restaurant-menu-updated",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildOutboxTopicNamefromEventType(tt.args.eventType); got != tt.want {
				t.Errorf("buildOutboxTopicNamefromEventType() = %v, want %v", got, tt.want)
			}
		})
	}
}

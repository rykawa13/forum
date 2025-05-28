package grpc

import (
	"context"
	"forum-service/internal/entity"
	"forum-service/internal/usecase"
	pb "forum-service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedForumServiceServer
	postUC usecase.PostUseCase
}

func NewServer(postUC usecase.PostUseCase) *Server {
	return &Server{
		postUC: postUC,
	}
}

func (s *Server) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.Post, error) {
	post, err := s.postUC.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "post not found")
	}

	return &pb.Post{
		Id:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		UserId:  int32(post.AuthorID),
		Author: &pb.User{
			Id:       int32(post.AuthorID),
			Username: post.Author.Username,
		},
		CreatedAt: timestamppb.New(post.CreatedAt),
		UpdatedAt: timestamppb.New(post.UpdatedAt),
	}, nil
}

func (s *Server) ListPosts(ctx context.Context, req *pb.ListPostsRequest) (*pb.ListPostsResponse, error) {
	posts, total, err := s.postUC.GetAll(ctx, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list posts")
	}

	pbPosts := make([]*pb.Post, len(posts))
	for i, post := range posts {
		pbPosts[i] = &pb.Post{
			Id:      post.ID,
			Title:   post.Title,
			Content: post.Content,
			UserId:  int32(post.AuthorID),
			Author: &pb.User{
				Id:       int32(post.AuthorID),
				Username: post.Author.Username,
			},
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		}
	}

	return &pb.ListPostsResponse{
		Posts: pbPosts,
		Total: int32(total),
	}, nil
}

func (s *Server) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	input := entity.CreatePostInput{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: int64(req.UserId),
	}

	post, err := s.postUC.Create(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create post")
	}

	return &pb.CreatePostResponse{
		Id: post.ID,
	}, nil
}

func (s *Server) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	input := entity.UpdatePostInput{
		Title:   req.Title,
		Content: req.Content,
	}

	err := s.postUC.Update(ctx, req.Id, input)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update post")
	}

	// Fetch the updated post
	post, err := s.postUC.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get updated post")
	}

	return &pb.Post{
		Id:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserId:    int32(post.AuthorID),
		CreatedAt: timestamppb.New(post.CreatedAt),
		UpdatedAt: timestamppb.New(post.UpdatedAt),
	}, nil
}

func (s *Server) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	err := s.postUC.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete post")
	}

	return &pb.DeletePostResponse{}, nil
}

func (s *Server) ListReplies(ctx context.Context, req *pb.ListRepliesRequest) (*pb.ListRepliesResponse, error) {
	replies, err := s.postUC.GetReplies(ctx, req.PostId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get replies")
	}

	pbReplies := make([]*pb.Reply, len(replies))
	for i, reply := range replies {
		pbReplies[i] = &pb.Reply{
			Id:      reply.ID,
			PostId:  reply.PostID,
			Content: reply.Content,
			UserId:  int32(reply.AuthorID),
			Author: &pb.User{
				Id:       int32(reply.AuthorID),
				Username: reply.Author.Username,
			},
			CreatedAt: timestamppb.New(reply.CreatedAt),
			UpdatedAt: timestamppb.New(reply.UpdatedAt),
		}
	}

	return &pb.ListRepliesResponse{
		Replies: pbReplies,
		Total:   int32(len(replies)),
	}, nil
}

func (s *Server) CreateReply(ctx context.Context, req *pb.CreateReplyRequest) (*pb.CreateReplyResponse, error) {
	input := entity.CreateReplyInput{
		Content:  req.Content,
		AuthorID: int64(req.UserId),
	}

	reply, err := s.postUC.CreateReply(ctx, req.PostId, input)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create reply")
	}

	return &pb.CreateReplyResponse{
		Id: reply.ID,
	}, nil
}

func (s *Server) DeleteReply(ctx context.Context, req *pb.DeleteReplyRequest) (*pb.DeleteReplyResponse, error) {
	err := s.postUC.DeleteReply(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete reply")
	}

	return &pb.DeleteReplyResponse{}, nil
}

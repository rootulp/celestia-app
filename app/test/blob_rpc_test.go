package app_test

// func TestBlobRPCQueries(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping blob RPC queries test in short mode.")
// 	}
// 	_, cctx := testnode.DefaultNetwork(t, time.Millisecond)
// 	h, err := cctx.WaitForHeightWithTimeout(10, time.Minute)
// 	require.NoError(t, err)
// 	require.Greater(t, h, int64(10))
// 	queryClient := types.NewQueryClient(cctx.GRPCClient)

// 	type test struct {
// 		name string
// 		req  func() error
// 	}
// 	tests := []test{
// 		{
// 			name: "blob by hash",
// 			req: func() error {
// 				_, err := queryClient.Params(
// 					context.Background(),
// 					// &types.QueryAttestationRequestByNonceRequest{Nonce: 1},
// 				)
// 				return err
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.req()
// 			assert.NoError(t, err)
// 		})
// 	}
// }

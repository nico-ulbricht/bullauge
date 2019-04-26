# bullauge
Tooling to expose the status of your Kubernetes cluster through GraphQL.  
A window into the boat. A bulls eye.

## installation
Installation can be done easily through the provided Helm chart.  
The attached ClusterRoleBinding will provide bullauge with read access on the entire cluster that it is hosted on.

```
helm install bullauge ./chart
```

For more details on configuration, please check https://github.com/nico-ulbricht/bullauge/blob/master/chart/values.yaml.

## schema
```graphql
schema {
  query: RootQuery
}

type POD {
  image: String
  logs(limit: Int = 10): [String!]!
  name: String
  namespace: String
  status: String
}

type RootQuery {
  pods(
    app: String
    namespace: String!
  ): [POD!]!
}
```